import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';
import '../../../../core/constants/colors.dart';
import '../../../../core/constants/tokens.dart';
import '../../../../core/constants/typography.dart';
import '../../../../core/widgets/neobrutal_card.dart';

/// Dark Neobrutalist Throughput Sparkline Chart.
class TelemetryChart extends StatelessWidget {
  final double currentVelocity;

  const TelemetryChart({super.key, required this.currentVelocity});

  @override
  Widget build(BuildContext context) {
    return NeobrutalCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text('THROUGHPUT STREAM', style: AegisTypography.cardLabel),
              Text(
                '60 FPS REALTIME',
                style: AegisTypography.actionButton.copyWith(
                  color: AegisColors.electricCyan,
                  fontSize: 9,
                ),
              ),
            ],
          ),
          const SizedBox(height: AegisTokens.space16),
          SizedBox(
            height: 90,
            child: LineChart(
              LineChartData(
                gridData: const FlGridData(show: false),
                titlesData: const FlTitlesData(show: false),
                borderData: FlBorderData(show: false),
                minX: 0,
                maxX: 6,
                minY: 0,
                maxY: 60,
                lineBarsData: [
                  LineChartBarData(
                    spots: [
                      const FlSpot(0, 12),
                      const FlSpot(1, 20),
                      const FlSpot(2, 14),
                      const FlSpot(3, 32),
                      const FlSpot(4, 25),
                      const FlSpot(5, 42),
                      FlSpot(6, currentVelocity > 0 ? currentVelocity : 35),
                    ],
                    isCurved: true,
                    curveSmoothness: 0.35,
                    color: AegisColors.electricCyan,
                    barWidth: 3.0,
                    isStrokeCapRound: true,
                    dotData: const FlDotData(show: false),
                    belowBarData: BarAreaData(
                      show: true,
                      color: AegisColors.electricCyan.withValues(alpha: 0.12),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}
