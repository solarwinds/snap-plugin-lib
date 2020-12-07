using Microsoft.Extensions.DependencyInjection;

namespace SnapPluginLib
{
    public static class ServiceCollectionExtensions
    {
        public static IServiceCollection AddSnapPlugin(this IServiceCollection services)
        {
            return services.AddSingleton<IRunner, Runner>();
        }
    }
}