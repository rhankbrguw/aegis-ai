import 'dart:async';
import 'dart:convert';
import 'package:web_socket_channel/web_socket_channel.dart';
import '../constants/config.dart';

/// Telemetry payload received from Go Gateway over WebSocket.
class TelemetryStreamPayload {
  final String model;
  final int totalTokens;
  final double estimatedCostUSD;
  final double dollarsSavedUSD;
  final int latencyMs;
  final bool isCacheHit;
  final String provider;
  final String circuitState;

  const TelemetryStreamPayload({
    required this.model,
    required this.totalTokens,
    required this.estimatedCostUSD,
    required this.dollarsSavedUSD,
    required this.latencyMs,
    required this.isCacheHit,
    required this.provider,
    required this.circuitState,
  });

  factory TelemetryStreamPayload.fromJson(Map<String, dynamic> json) {
    return TelemetryStreamPayload(
      model: json['model'] as String? ?? 'unknown',
      totalTokens: (json['total_tokens'] as num?)?.toInt() ?? 0,
      estimatedCostUSD: (json['estimated_cost_usd'] as num?)?.toDouble() ?? 0.0,
      dollarsSavedUSD: (json['dollars_saved_usd'] as num?)?.toDouble() ?? 0.0,
      latencyMs: (json['latency_ms'] as num?)?.toInt() ?? 0,
      isCacheHit: json['is_cache_hit'] as bool? ?? false,
      provider: json['provider'] as String? ?? 'primary',
      circuitState: json['circuit_state'] as String? ?? 'CLOSED',
    );
  }
}

/// WebSocket telemetry client with resilient connection lifecycle.
class TelemetryWebSocketClient {
  final String url;
  WebSocketChannel? _channel;
  StreamController<TelemetryStreamPayload>? _controller;
  bool _isDisposed = false;

  TelemetryWebSocketClient({String? url}) : url = url ?? AegisConfig.wsTelemetryUrl;

  Stream<TelemetryStreamPayload> connect() {
    _controller = StreamController<TelemetryStreamPayload>.broadcast(
      onListen: _initiateConnection,
      onCancel: _closeConnection,
    );
    return _controller!.stream;
  }

  void _initiateConnection() {
    if (_isDisposed) return;
    try {
      _channel = WebSocketChannel.connect(Uri.parse(url));
      _channel!.stream.listen(
        _handleMessage,
        onError: (_) => _scheduleReconnect(),
        onDone: () => _scheduleReconnect(),
        cancelOnError: true,
      );
    } catch (_) {
      _scheduleReconnect();
    }
  }

  void _handleMessage(dynamic message) {
    try {
      final data = jsonDecode(message.toString()) as Map<String, dynamic>;
      final payload = TelemetryStreamPayload.fromJson(data);
      if (_controller != null && !_controller!.isClosed) {
        _controller!.add(payload);
      }
    } catch (_) {}
  }

  void _scheduleReconnect() {
    if (_isDisposed) return;
    Timer(const Duration(seconds: 3), _initiateConnection);
  }

  void _closeConnection() {
    _channel?.sink.close();
    _channel = null;
  }

  void dispose() {
    _isDisposed = true;
    _closeConnection();
    _controller?.close();
  }
}
