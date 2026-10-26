import 'package:flutter/material.dart';
import '../../../../core/constants/colors.dart';
import '../../../../core/constants/tokens.dart';
import '../../../../core/constants/typography.dart';

/// Industrial Hardware & Edge Status Badge Strip.
class StatusStrip extends StatelessWidget {
  const StatusStrip({super.key});

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      child: Row(
        children: [
          _buildBadge('TLS 1.3 // EDGE', AegisColors.electricLime),
          const SizedBox(width: AegisTokens.space8),
          _buildBadge('REGION: SG-SIN', AegisColors.electricCyan),
          const SizedBox(width: AegisTokens.space8),
          _buildBadge('CACHE: SUB-15MS', AegisColors.cyberYellow),
          const SizedBox(width: AegisTokens.space8),
          _buildBadge('UPTIME: 99.99%', AegisColors.textSecondary),
        ],
      ),
    );
  }

  Widget _buildBadge(String label, Color accent) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: AegisColors.surfaceSecondary,
        borderRadius: BorderRadius.circular(4),
        border: Border.all(color: AegisColors.surfaceBorder, width: 1.2),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            width: 5,
            height: 5,
            decoration: BoxDecoration(color: accent, shape: BoxShape.circle),
          ),
          const SizedBox(width: 6),
          Text(
            label,
            style: AegisTypography.actionButton.copyWith(
              color: AegisColors.textSecondary,
              fontSize: 9,
              letterSpacing: 1.0,
            ),
          ),
        ],
      ),
    );
  }
}
