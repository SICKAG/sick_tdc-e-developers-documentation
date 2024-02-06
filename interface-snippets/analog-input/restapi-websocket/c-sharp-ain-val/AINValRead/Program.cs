using System;
using System.IO;
using System.Net.WebSockets;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using System.Collections.Generic;
using System.Net.Http;
using System.Text.Json;

public class AnalogValueChange
{
    public string AinName { get; set; }
    public double PreviousValue { get; set; }
    public double NewValue { get; set; }
}

public class Program
{

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


    public static void Main(string[] args)
    {
        var tasks = new List<Task>();

        tasks.Add(Task.Run(() => ListenAinVal()));

        Task.WaitAll(tasks.ToArray());
    }
}
