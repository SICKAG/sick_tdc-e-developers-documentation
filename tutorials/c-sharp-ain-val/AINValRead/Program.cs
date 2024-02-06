using System;
using System.IO;
using System.Net.WebSockets;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using System.Collections.Generic;
using System.Net.Http;
using System.Text.Json;

public class Analog
{
    public string AinName { get; set; }
    public string State { get; set; }
}

public class AnalogVal
{
    public string AinName { get; set; }
    public double Value { get; set; }
}

public class AnalogValueChange
{
    public string AinName { get; set; }
    public double PreviousValue { get; set; }
    public double NewValue { get; set; }
}

public class TokenResponse
{
    public string Token { get; set; }
}

public class Program
{
    private static string token;

    private static async Task<string> GetTokenAsync()
    {
        var tokenURL = "http://192.168.0.100:59801/user/Service/token";
        var password = "servicelevel";

        var formContent = new FormUrlEncodedContent(new Dictionary<string, string>
        {
            { "password", password }
        });

        using (var client = new HttpClient())
        {
            var response = await client.PostAsync(tokenURL, formContent);
            var json = await response.Content.ReadAsStringAsync();
            var tokenResp = JsonSerializer.Deserialize<TokenResponse>(json);
            return tokenResp.Token;
        }
    }

    private static async Task<List<Analog>> GetAllAnalogStatesAsync(string token)
    {
        var url = "http://192.168.0.100:59801/tdce/analog-inputs/GetStates";

        using (var client = new HttpClient())
        {
            client.DefaultRequestHeaders.Add("Authorization", "Bearer " + token);

            var json = await client.GetStringAsync(url);
            return JsonSerializer.Deserialize<List<Analog>>(json);
        }
    }

    private static async Task<List<AnalogVal>> GetAllAnalogValuesAsync(string token)
    {
        var url = "http://192.168.0.100:59801/tdce/analog-inputs/GetValues";

        using (var client = new HttpClient())
        {
            client.DefaultRequestHeaders.Add("Authorization", "Bearer " + token);

            var json = await client.GetStringAsync(url);
            return JsonSerializer.Deserialize<List<AnalogVal>>(json);
        }
    }

    private static void PrintAnalogState(string ainName, string state)
    {
        Console.WriteLine($"\nAnalog Input {ainName}, State: {state}");
    }

    private static async Task<double> GetAnalogValueAsync(string token, string ain)
    {
        var url = $"http://192.168.0.100:59801/tdce/analog-inputs/GetValue/{ain}";

        using (var client = new HttpClient())
        {
            client.DefaultRequestHeaders.Add("Authorization", "Bearer " + token);

            var response = await client.GetAsync(url);
            var valueStr = await response.Content.ReadAsStringAsync();
            return double.Parse(valueStr);
        }
    }

    private static void ListenAinVal()
    {
        var serverUrl = new Uri("ws://192.168.0.100:31768/ws/tdce/analog-inputs/value");
        using (var client = new ClientWebSocket())
        {
            client.ConnectAsync(serverUrl, CancellationToken.None).Wait();
            Console.WriteLine("Connected to WebSocket");

            var receiveBuffer = new ArraySegment<byte>(new byte[1024]);

            while (true)
            {
                var result = client.ReceiveAsync(receiveBuffer, CancellationToken.None).Result;
                var message = Encoding.UTF8.GetString(receiveBuffer.Array, 0, result.Count);
                var avChange = JsonSerializer.Deserialize<AnalogValueChange>(message);

                Console.WriteLine($"Received AnalogValueChange: AinName={avChange.AinName}, PreviousValue={avChange.PreviousValue}, NewValue={avChange.NewValue}");
            }
        }
    }

    private static async Task HttpGetAinAsync()
    {
        token = await GetTokenAsync();

        var analogs = await GetAllAnalogStatesAsync(token);
        Console.WriteLine("Printing all analog states:");
        foreach (var analog in analogs)
        {
            Console.WriteLine($"AinName: {analog.AinName}, State: {analog.State}");
        }

        var analogVals = await GetAllAnalogValuesAsync(token);
        Console.WriteLine("\nAnalog Values:");
        foreach (var analogVal in analogVals)
        {
            Console.WriteLine($"Analog Name: {analogVal.AinName}, Value: {analogVal.Value}");
        }

        var ain = "AIN_A";
        PrintAnalogState(ain, "AIN_A");
        var value = await GetAnalogValueAsync(token, ain);
        Console.WriteLine($"\nAnalog Name: {ain}, Value: {value}");

        await Task.Delay(TimeSpan.FromSeconds(2));
    }

    public static void Main(string[] args)
    {
        var tasks = new List<Task>();

        tasks.Add(HttpGetAinAsync());
        tasks.Add(Task.Run(() => ListenAinVal()));

        Task.WaitAll(tasks.ToArray());
    }
}
