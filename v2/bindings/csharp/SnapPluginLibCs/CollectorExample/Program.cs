using System;
using Autofac;
using Autofac.Extensions.DependencyInjection;
using Microsoft.Extensions.DependencyInjection;
using SnapPluginLib;

namespace CollectorExample
{
    class Program
    {
        static void Main(string[] args)
        {
            using (var scope = BuildContainer().BeginLifetimeScope())
            {
                var runner = scope.Resolve<IRunner>();
                var collectorPlugin = scope.Resolve<ICollectorPlugin>();

                runner.StartCollector(collectorPlugin);
            }
        }

        private static IContainer BuildContainer()
        {
            var containerBuilder = new ContainerBuilder();
            var serviceCollection = new ServiceCollection();
            serviceCollection
                .AddSnapPlugin()
                .AddSingleton<ICollectorPlugin>(provider => new CollectorExample("collector-example", new Version(0, 0, 1)));
           
            containerBuilder.Populate(serviceCollection);

            return containerBuilder.Build();
        }
    }
}