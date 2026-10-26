import 'package:flutter/material.dart';
import '../../../../core/constants/colors.dart';
import '../../../../core/constants/tokens.dart';
import '../../../../core/constants/typography.dart';
import '../../../../core/widgets/interactive_kill_slider.dart';
import '../../../../core/widgets/neobrutal_card.dart';

/// Dark Neobrutalist Tactical Emergency Control.
class KillSwitchControl extends StatelessWidget {
  final bool isKillSwitchActive;
  final ValueChanged<bool> onToggle;
  final VoidCallback onSimulateOutage;
  final VoidCallback onRecoverCircuit;

  const KillSwitchControl({
    super.key,
    required this.isKillSwitchActive,
    required this.onToggle,
    required this.onSimulateOutage,
    required this.onRecoverCircuit,
  });

  @override
  Widget build(BuildContext context) {
    return NeobrutalCard(
      borderColor: isKillSwitchActive ? AegisColors.laserCrimson : AegisColors.surfaceBorder,
      shadowColor: isKillSwitchActive ? AegisColors.laserCrimson.withValues(alpha: 0.3) : AegisColors.hardBlack,
      backgroundColor: isKillSwitchActive ? AegisColors.laserCrimson.withValues(alpha: 0.08) : AegisColors.surfacePrimary,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          InteractiveKillSlider(
            isKillSwitchActive: isKillSwitchActive,
            onToggle: onToggle,
          ),
          const SizedBox(height: AegisTokens.space12),
          Row(
            children: [
              Expanded(
                child: _buildNeobrutalButton(
                  label: 'INJECT 429 OUTAGE',
                  icon: Icons.bolt,
                  color: AegisColors.cyberYellow,
                  onTap: onSimulateOutage,
                ),
              ),
              const SizedBox(width: AegisTokens.space12),
              Expanded(
                child: _buildNeobrutalButton(
                  label: 'RESET BREAKER',
                  icon: Icons.restart_alt,
                  color: AegisColors.electricLime,
                  onTap: onRecoverCircuit,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildNeobrutalButton({
    required String label,
    required IconData icon,
    required Color color,
    required VoidCallback onTap,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 10),
        decoration: BoxDecoration(
          color: AegisColors.surfaceSecondary,
          borderRadius: BorderRadius.circular(AegisTokens.radiusSmall),
          border: Border.all(color: color.withValues(alpha: 0.6), width: 1.5),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(icon, color: color, size: 14),
            const SizedBox(width: 6),
            Text(
              label,
              style: AegisTypography.actionButton.copyWith(color: color, fontSize: 10),
            ),
          ],
        ),
      ),
    );
  }
}
