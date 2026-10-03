import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/constants/colors.dart';
import '../../../core/constants/tokens.dart';
import '../../../core/constants/typography.dart';
import '../../../core/widgets/cyber_grid_background.dart';
import 'dashboard_controller.dart';
import 'widgets/header_bar.dart';
import 'widgets/kill_switch_control.dart';
import 'widgets/metrics_grid.dart';
import 'widgets/status_strip.dart';
import 'widgets/telemetry_chart.dart';
import 'widgets/terminal_logs.dart';

/// Dark Neobrutalism SRE Dashboard View with Cyber Grid & Fluid Motion.
class DashboardView extends ConsumerWidget {
  const DashboardView({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(dashboardProvider);
    final notifier = ref.read(dashboardProvider.notifier);

    return Scaffold(
      backgroundColor: AegisColors.background,
      appBar: _buildAppBar(),
      body: CyberGridBackground(
        child: SingleChildScrollView(
          padding: const EdgeInsets.symmetric(
            horizontal: AegisTokens.space16,
            vertical: AegisTokens.space8,
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: _buildAnimatedSections(state, notifier),
          ),
        ),
      ),
    );
  }

  PreferredSizeWidget _buildAppBar() {
    return AppBar(
      title: Text('AEGIS // SENTINEL CORE', style: AegisTypography.heroTitle),
      centerTitle: false,
      backgroundColor: AegisColors.background,
      elevation: 0,
    );
  }

  List<Widget> _buildAnimatedSections(DashboardState state, DashboardNotifier notifier) {
    return [
      TacticalHeaderBar(
        circuitState: state.circuitState,
        isKillSwitchActive: state.isKillSwitchActive,
      ).animate().fadeIn(duration: 350.ms).slideY(begin: 0.08, end: 0),
      const SizedBox(height: AegisTokens.space8),
      const StatusStrip().animate().fadeIn(delay: 50.ms, duration: 350.ms),
      const SizedBox(height: AegisTokens.space12),
      MetricsGrid(
        tokensPerSec: state.tokensPerSec,
        totalSavedUSD: state.totalSavedUSD,
        circuitState: state.circuitState,
      ).animate().fadeIn(delay: 100.ms, duration: 350.ms).slideY(begin: 0.08, end: 0),
      const SizedBox(height: AegisTokens.space12),
      TelemetryChart(currentVelocity: state.tokensPerSec)
          .animate().fadeIn(delay: 180.ms, duration: 350.ms).slideY(begin: 0.08, end: 0),
      const SizedBox(height: AegisTokens.space12),
      TerminalLogsView(logs: state.recentEvents)
          .animate().fadeIn(delay: 260.ms, duration: 350.ms).slideY(begin: 0.08, end: 0),
      const SizedBox(height: AegisTokens.space12),
      KillSwitchControl(
        isKillSwitchActive: state.isKillSwitchActive,
        onToggle: (active) => notifier.toggleKillSwitch(active),
        onSimulateOutage: () => _handleSimulateOutage(notifier),
        onRecoverCircuit: () => _handleRecoverCircuit(notifier),
      ).animate().fadeIn(delay: 340.ms, duration: 350.ms).slideY(begin: 0.08, end: 0),
      const SizedBox(height: AegisTokens.space24),
    ];
  }

  void _handleSimulateOutage(DashboardNotifier notifier) {
    notifier.setCircuitState('OPEN');
    notifier.updateMetrics(tps: 48.5, upm: 0.15, saved: 0.082);
    final time = DateTime.now().toIso8601String().substring(11, 19);
    notifier.addEvent('[$time] [SIMULATED_429] -> Tripped Circuit to OPEN');
  }

  void _handleRecoverCircuit(DashboardNotifier notifier) {
    notifier.setCircuitState('CLOSED');
    notifier.updateMetrics(tps: 18.0, upm: 0.03, saved: 0.015);
    final time = DateTime.now().toIso8601String().substring(11, 19);
    notifier.addEvent('[$time] [RESET_BREAKER] -> Restored state to CLOSED');
  }
}
