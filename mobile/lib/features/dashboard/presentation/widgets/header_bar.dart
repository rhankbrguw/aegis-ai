import 'package:flutter/material.dart';
import '../../../../core/constants/colors.dart';
import '../../../../core/constants/tokens.dart';
import '../../../../core/constants/typography.dart';
import '../../../../core/widgets/neobrutal_card.dart';
import '../../../../core/widgets/pulsing_beacon.dart';

/// Dark Neobrutalist Header Bar with Live Pulsing Beacon.
class TacticalHeaderBar extends StatelessWidget {
  final String circuitState;
  final bool isKillSwitchActive;

  const TacticalHeaderBar({
    super.key,
    required this.circuitState,
    required this.isKillSwitchActive,
  });

  @override
  Widget build(BuildContext context) {
    final isOnline = circuitState == 'CLOSED' && !isKillSwitchActive;
    final statusColor = isOnline ? AegisColors.electricLime : AegisColors.laserCrimson;
    final statusLabel = isKillSwitchActive ? 'HALTED' : (circuitState == 'OPEN' ? 'FAILOVER' : 'OPERATIONAL');

    return NeobrutalCard(
      borderColor: statusColor.withValues(alpha: 0.5),
      shadowColor: statusColor.withValues(alpha: 0.2),
      padding: const EdgeInsets.symmetric(horizontal: AegisTokens.space16, vertical: AegisTokens.space12),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Row(
            children: [
              PulsingBeacon(color: statusColor, size: 8),
              const SizedBox(width: AegisTokens.space12),
              Text(
                'api.rhankbrguw.xyz',
                style: AegisTypography.sectionTitle.copyWith(fontSize: 12),
              ),
            ],
          ),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
            decoration: BoxDecoration(
              color: statusColor,
              borderRadius: BorderRadius.circular(4),
              border: Border.all(color: AegisColors.hardBlack, width: 1.5),
            ),
            child: Text(
              statusLabel,
              style: AegisTypography.microBadge,
            ),
          ),
        ],
      ),
    );
  }
}
