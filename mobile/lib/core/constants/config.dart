/// Network and Gateway Target Configuration.
class AegisConfig {
  const AegisConfig._();

  static const String appEnv = String.fromEnvironment('APP_ENV', defaultValue: 'development');
  static const bool isProduction = appEnv == 'production';

  static const String localHost = '127.0.0.1';
  static const int defaultPort = 8080;
  static const String prodDomain = 'api-aegis.rhankbrguw.xyz';

  static const String httpBaseUrl = isProduction
      ? 'https://api-aegis.rhankbrguw.xyz'
      : 'http://127.0.0.1:8080';

  static const String wsTelemetryUrl = isProduction
      ? 'wss://api-aegis.rhankbrguw.xyz/v1/telemetry/stream'
      : 'ws://127.0.0.1:8080/v1/telemetry/stream';

  static const String circuitOverrideUrl = isProduction
      ? 'https://api-aegis.rhankbrguw.xyz/v1/circuit/override'
      : 'http://127.0.0.1:8080/v1/circuit/override';

  static const String defaultApiKey = 'aegis_live_83734d2232945f4899da40edcc46d75a216081a08266be30';
}


