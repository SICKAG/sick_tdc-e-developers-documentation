using DioLED;
using System;
using System.Diagnostics;
using System.Net;
using System.Net.Http;
using System.Text;
using System.Threading.Tasks;

Stopwatch stopwatch = Stopwatch.StartNew();
bool status;
string token;
bool started = false;
HttpClient client = new HttpClient();
long totalGlowingTime = 0;

//object locks; semaphores
object ledLock = new object();
object totalTimeLock = new object();
SemaphoreSlim ledSemaphore = new SemaphoreSlim(1);

using (client)
{
    await SetupAndRun();
}

//sets up environment, then enters loop
async Task SetupAndRun()
{
    string tokenJson = await GetTokenAsync(client);
    token = SplitToken(tokenJson);
    client.DefaultRequestHeaders.Add("Authorization", "Bearer " + token);

    //creating DIO js
    var endpoint = new Uri("http://192.168.0.100:59801/tdce/dio/SetStates");
    await OnRepeatAsync(client, endpoint);
}

//infinite loop; works with stopwatches, http requests and database
async Task OnRepeatAsync(HttpClient client, Uri endpoint)
{
    while (true)
    {
        status = await FetchStatusAsync(client);
        if (!status)
        {
            if (started)
                StopStopWatch(endpoint);
        }
        else
        {
            if (!started)
                StartStopwatch();
        }
    }
}

//function starts the stopwatch counter
void StartStopwatch()
{
    stopwatch.Start();
    started = true;
}

//function stops stopwatch counter, sets total time and starts a new task to status
void StopStopWatch(Uri endpoint)
{
    started = false;
    stopwatch.Stop();

    long elapsedMilliseconds;
    //locks to ensure concurrency
    lock (totalTimeLock)
    {
        elapsedMilliseconds = stopwatch.ElapsedMilliseconds;
        totalGlowingTime += elapsedMilliseconds;
    }
    stopwatch.Reset();

    string js = "[{ \"DioName\": \"DIO_A\", \"Value\": 1, \"Direction\": \"Output\" }]";
    //makes a new task to run simultaneously
    Task.Factory.StartNew(() => SetStatusAsync(js, client, endpoint, elapsedMilliseconds));
}

//task sets status and controls LED
async Task SetStatusAsync(string js, HttpClient client, Uri endpoint, long elapsedMilliseconds)
{
    //set first status
    var pay = new StringContent(js, Encoding.UTF8, "application/json");
    var resp = await client.PostAsync(endpoint, pay);
    Console.WriteLine("Total Glowing Time: " + totalGlowingTime);

    //if value already larger than 0
    if (totalGlowingTime > 0)
    {
        await ledSemaphore.WaitAsync();
        try
        {
            //set LED on
            var payload = new StringContent(js, Encoding.UTF8, "application/json");
            var response = await client.PostAsync(endpoint, payload);

            //calculating remaining time
            DateTime startGlowTime = DateTime.Now;
            while (totalGlowingTime > 0)
            {
                long remainingTime = totalGlowingTime - (long)(DateTime.Now - startGlowTime).TotalMilliseconds;
                if (remainingTime <= 0)
                {
                    break;
                }
                //minimum delay for checking if updated
                await Task.Delay(1);
            }

            //set LED off
            js = "[{ \"DioName\": \"DIO_A\", \"Value\": 0, \"Direction\": \"Output\" }]";
            payload = new StringContent(js, Encoding.UTF8, "application/json");
            response = await client.PostAsync(endpoint, payload);

            if (totalGlowingTime > 0)
            {
                DatabaseHandler dbHandler = new DatabaseHandler(totalGlowingTime.ToString());
                dbHandler.ConnectToDB();
                Console.WriteLine("Sleep: " + totalGlowingTime);
                totalGlowingTime = 0;
            }
        }
        finally
        {
            ledSemaphore.Release();
        }
    }
}

//function reformats token appearance
string SplitToken(string tokenJson)
{
    string[] split = tokenJson.Split(':');
    string splat = split[1];
    //removes last two characters
    return splat.Substring(1, splat.Length - 3);
}

//function fetches authentication token
async Task<string> GetTokenAsync(HttpClient client)
{
    var endp = new Uri("http://192.168.0.100:59801/user/Service/token");
    var dict = new Dictionary<string, string>();
    dict.Add("password", "servicelevel");
    var content = new FormUrlEncodedContent(dict);
    var response = await client.PostAsync(endp, content);
    return await response.Content.ReadAsStringAsync();
}

//function fetches status of the DIO state B
async Task<bool> FetchStatusAsync(HttpClient client)
{
    var endpoint = new Uri("http://192.168.0.100:59801/tdce/dio/GetState/DIO_B");
    var jsonResult = await client.GetStringAsync(endpoint);
    var val = jsonResult.Split(',');
    val = val[2].Split(':');
    return val[1] == "1";
}
