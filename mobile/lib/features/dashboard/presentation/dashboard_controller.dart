import 'dart:async';
import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/constants/config.dart';
import '../../../core/network/websocket_client.dart';

/// Telemetry state snapshot for the dashboard.
class DashboardState {
  final double tokensPerSec;
  final double usdPerMin;
  final double totalSavedUSD;
  final String circuitState;
  final bool isKillSwitchActive;
  final List<String> recentEvents;

  const DashboardState({
    this.tokensPerSec = 0.0,
    this.usdPerMin = 0.0,
    this.totalSavedUSD = 0.0,
    this.circuitState = 'CLOSED',
    this.isKillSwitchActive = false,
    this.recentEvents = const [],
  });

  DashboardState copyWith({
    double? tokensPerSec,
    double? usdPerMin,
    double? totalSavedUSD,
    String? circuitState,
    bool? isKillSwitchActive,
    List<String>? recentEvents,
  }) {
    return DashboardState(
      tokensPerSec: tokensPerSec ?? this.tokensPerSec,
      usdPerMin: usdPerMin ?? this.usdPerMin,
      totalSavedUSD: totalSavedUSD ?? this.totalSavedUSD,
      circuitState: circuitState ?? this.circuitState,
      isKillSwitchActive: isKillSwitchActive ?? this.isKillSwitchActive,
      recentEvents: recentEvents ?? this.recentEvents,
    );
  }
}

/// Controller managing telemetry updates, remote kill-switch, and WS stream.
class DashboardNotifier extends StateNotifier<DashboardState> {
  final TelemetryWebSocketClient? _wsClient;
  final Dio _dio;
  StreamSubscription<TelemetryStreamPayload>? _wsSubscription;

  DashboardNotifier({TelemetryWebSocketClient? wsClient, Dio? dio})
      : _wsClient = wsClient,
        _dio = dio ?? Dio(),
        super(const DashboardState()) {
    _subscribeToLiveTelemetry();
  }

  void _subscribeToLiveTelemetry() {
    if (_wsClient == null) return;
    _wsSubscription = _wsClient!.connect().listen((payload) {
      final tps = payload.latencyMs > 0 ? (payload.totalTokens / (payload.latencyMs / 1000.0)) : 0.0;
      final timeStr = DateTime.now().toIso8601String().substring(11, 19);
      final cacheTag = payload.isCacheHit ? '[CACHE_HIT]' : '[${payload.provider.toUpperCase()}]';
      final logStr = '[$timeStr] ${payload.model} -> $cacheTag ${payload.latencyMs}ms (${payload.totalTokens} toks)';

      state = state.copyWith(
        tokensPerSec: double.parse(tps.toStringAsFixed(1)),
        totalSavedUSD: state.totalSavedUSD + payload.dollarsSavedUSD,
        circuitState: payload.circuitState,
      );
      addEvent(logStr);
    });
  }

  void updateMetrics({required double tps, required double upm, required double saved}) {
    state = state.copyWith(
      tokensPerSec: tps,
      usdPerMin: upm,
      totalSavedUSD: state.totalSavedUSD + saved,
    );
  }

  void setCircuitState(String newState) {
    state = state.copyWith(circuitState: newState);
  }

  Future<void> toggleKillSwitch(bool active) async {
    state = state.copyWith(
      isKillSwitchActive: active,
      circuitState: active ? 'OPEN' : 'CLOSED',
    );
    try {
      await _dio.post(
        AegisConfig.circuitOverrideUrl,
        data: {'kill_switch': active},
        options: Options(sendTimeout: const Duration(seconds: 2)),
      );
    } catch (_) {}
  }

  void addEvent(String eventLog) {
    final updated = [eventLog, ...state.recentEvents];
    if (updated.length > 20) updated.removeLast();
    state = state.copyWith(recentEvents: updated);
  }

  @override
  void dispose() {
    _wsSubscription?.cancel();
    _wsClient?.dispose();
    super.dispose();
  }
}

final telemetryClientProvider = Provider<TelemetryWebSocketClient>((ref) {
  return TelemetryWebSocketClient();
});

final dashboardProvider = StateNotifierProvider<DashboardNotifier, DashboardState>((ref) {
  final wsClient = ref.watch(telemetryClientProvider);
  return DashboardNotifier(wsClient: wsClient);
});
