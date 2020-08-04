using System;
using System.Text.RegularExpressions;
using SnapPluginLib;

namespace CollectorExample
{
    class Program
    {
        static void Main(string[] args)
        {
            Runner.StartCollector(new CollectorExample("collector-example", "v0.0.1"));
        }
    }
}
