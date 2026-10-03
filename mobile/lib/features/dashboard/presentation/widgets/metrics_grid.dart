import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import '../../../../core/constants/colors.dart';
import '../../../../core/constants/tokens.dart';
import '../../../../core/constants/typography.dart';
import '../../../../core/widgets/animated_counter.dart';
import '../../../../core/widgets/neobrutal_card.dart';

/// Dark Neobrutalist SRE Metrics Grid with Tactile Motion & Typography.
class MetricsGrid extends StatelessWidget {
  final double tokensPerSec;
  final double totalSavedUSD;
  final String circuitState;

  const MetricsGrid({
    super.key,
    required this.tokensPerSec,
    required this.totalSavedUSD,
    required this.circuitState,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Expanded(child: _buildVelocityCard()),
        const SizedBox(width: AegisTokens.space12),
        Expanded(child: _buildSavingsCard()),
      ],
    );
  }

  Widget _buildVelocityCard() {
    return NeobrutalCard(
      borderColor: AegisColors.electricCyan.withValues(alpha: 0.4),
      shadowColor: AegisColors.electricCyan.withValues(alpha: 0.15),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text('BURN VELOCITY', style: AegisTypography.cardLabel),
              const Icon(Icons.speed, color: AegisColors.electricCyan, size: 16),
            ],
          ),
          const SizedBox(height: AegisTokens.space12),
          AnimatedCounter(
            value: tokensPerSec,
            suffix: ' T/s',
            decimalPlaces: 1,
            style: AegisTypography.telemetryMetric.copyWith(color: AegisColors.electricCyan),
          ),
        ],
      ),
    );
  }

  Widget _buildSavingsCard() {
    return NeobrutalCard(
      borderColor: AegisColors.electricLime.withValues(alpha: 0.4),
      shadowColor: AegisColors.electricLime.withValues(alpha: 0.15),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text('FINOPS SAVINGS', style: AegisTypography.cardLabel),
              const Icon(Icons.savings_outlined, color: AegisColors.electricLime, size: 16),
            ],
          ),
          const SizedBox(height: AegisTokens.space12),
          AnimatedCounter(
            value: totalSavedUSD,
            prefix: r'$',
            decimalPlaces: 4,
            style: AegisTypography.telemetryMetric,
          ),
        ],
      ),
    ).animate(onPlay: (c) => c.repeat(reverse: true))
     .shimmer(duration: 3500.ms, color: AegisColors.electricLime.withValues(alpha: 0.08));
  }
}
