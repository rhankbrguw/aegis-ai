import 'package:flutter/material.dart';
import '../constants/colors.dart';
import '../constants/tokens.dart';

/// Tactile Slide-to-Halt Emergency Action Slider with Spring Physics.
class InteractiveKillSlider extends StatefulWidget {
  final bool isKillSwitchActive;
  final ValueChanged<bool> onToggle;

  const InteractiveKillSlider({
    super.key,
    required this.isKillSwitchActive,
    required this.onToggle,
  });

  @override
  State<InteractiveKillSlider> createState() => _InteractiveKillSliderState();
}

class _InteractiveKillSliderState extends State<InteractiveKillSlider> {
  double _dragPosition = 0.0;
  static const double _thumbWidth = 44.0;

  @override
  Widget build(BuildContext context) {
    final activeColor = widget.isKillSwitchActive ? AegisColors.laserCrimson : AegisColors.electricLime;

    return LayoutBuilder(
      builder: (context, constraints) {
        final maxDrag = constraints.maxWidth - _thumbWidth - 8.0;

        return Container(
          height: 52,
          padding: const EdgeInsets.symmetric(horizontal: 4.0),
          decoration: BoxDecoration(
            color: widget.isKillSwitchActive
                ? AegisColors.laserCrimson.withValues(alpha: 0.12)
                : AegisColors.surfaceSecondary,
            borderRadius: BorderRadius.circular(AegisTokens.radiusMedium),
            border: Border.all(
              color: widget.isKillSwitchActive
                  ? AegisColors.laserCrimson
                  : AegisColors.surfaceBorder,
            ),
          ),
          child: Stack(
            alignment: Alignment.centerLeft,
            children: [
              Center(
                child: Text(
                  widget.isKillSwitchActive ? 'GATEWAY HALTED (TAP TO RESTORE)' : 'SLIDE TO EMERGENCY HALT >>>',
                  style: TextStyle(
                    color: widget.isKillSwitchActive ? AegisColors.laserCrimson : AegisColors.textMuted,
                    fontSize: 10,
                    fontWeight: FontWeight.bold,
                    letterSpacing: 1.2,
                    fontFamily: 'monospace',
                  ),
                ),
              ),
              Positioned(
                left: widget.isKillSwitchActive ? maxDrag : _dragPosition,
                child: GestureDetector(
                  onHorizontalDragUpdate: (details) {
                    if (widget.isKillSwitchActive) return;
                    setState(() {
                      _dragPosition = (_dragPosition + details.delta.dx).clamp(0.0, maxDrag);
                    });
                  },
                  onHorizontalDragEnd: (details) {
                    if (widget.isKillSwitchActive) return;
                    if (_dragPosition >= maxDrag * 0.75) {
                      widget.onToggle(true);
                    }
                    setState(() => _dragPosition = 0.0);
                  },
                  onTap: () {
                    if (widget.isKillSwitchActive) {
                      widget.onToggle(false);
                      setState(() => _dragPosition = 0.0);
                    }
                  },
                  child: Container(
                    width: _thumbWidth,
                    height: 44,
                    decoration: BoxDecoration(
                      color: activeColor,
                      borderRadius: BorderRadius.circular(AegisTokens.radiusSmall),
                      boxShadow: [
                        BoxShadow(
                          color: activeColor.withValues(alpha: 0.5),
                          blurRadius: 8,
                          spreadRadius: 1,
                        ),
                      ],
                    ),
                    child: Icon(
                      widget.isKillSwitchActive ? Icons.lock_open : Icons.power_settings_new,
                      color: AegisColors.background,
                      size: 20,
                    ),
                  ),
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}
