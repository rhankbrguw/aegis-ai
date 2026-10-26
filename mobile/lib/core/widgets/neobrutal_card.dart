import 'package:flutter/material.dart';
import '../constants/colors.dart';
import '../constants/tokens.dart';

/// Authentic Dark Neobrutalism Card with Hard Offset Shadow & Physical Press Physics.
class NeobrutalCard extends StatefulWidget {
  final Widget child;
  final EdgeInsetsGeometry? padding;
  final Color? borderColor;
  final Color? backgroundColor;
  final Color? shadowColor;
  final double shadowOffset;
  final VoidCallback? onTap;

  const NeobrutalCard({
    super.key,
    required this.child,
    this.padding = const EdgeInsets.all(AegisTokens.space16),
    this.borderColor,
    this.backgroundColor,
    this.shadowColor,
    this.shadowOffset = 3.5,
    this.onTap,
  });

  @override
  State<NeobrutalCard> createState() => _NeobrutalCardState();
}

class _NeobrutalCardState extends State<NeobrutalCard> {
  bool _isPressed = false;

  @override
  Widget build(BuildContext context) {
    final effectiveShadow = widget.shadowColor ?? AegisColors.hardBlack;
    final border = widget.borderColor ?? AegisColors.surfaceBorder;
    final currentOffset = _isPressed ? 1.0 : widget.shadowOffset;

    return GestureDetector(
      onTapDown: (_) => widget.onTap != null ? setState(() => _isPressed = true) : null,
      onTapUp: (_) => widget.onTap != null ? setState(() => _isPressed = false) : null,
      onTapCancel: () => widget.onTap != null ? setState(() => _isPressed = false) : null,
      onTap: widget.onTap,
      child: Transform.translate(
        offset: _isPressed ? const Offset(2.0, 2.0) : Offset.zero,
        child: Container(
          padding: widget.padding,
          decoration: BoxDecoration(
            color: widget.backgroundColor ?? AegisColors.surfacePrimary,
            borderRadius: BorderRadius.circular(AegisTokens.radiusSmall),
            border: Border.all(color: border, width: 2.0),
            boxShadow: [
              BoxShadow(
                color: effectiveShadow,
                offset: Offset(currentOffset, currentOffset),
                blurRadius: 0,
                spreadRadius: 0,
              ),
            ],
          ),
          child: widget.child,
        ),
      ),
    );
  }
}
