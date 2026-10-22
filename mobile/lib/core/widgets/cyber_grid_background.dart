import 'package:flutter/material.dart';
import '../constants/colors.dart';

/// Minimalist Dark Cyber Grid Background with Radial Ambient Glow.
class CyberGridBackground extends StatelessWidget {
  final Widget child;

  const CyberGridBackground({super.key, required this.child});

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        Positioned.fill(
          child: CustomPaint(
            painter: _CyberGridPainter(),
          ),
        ),
        Positioned.fill(
          child: Container(
            decoration: BoxDecoration(
              gradient: RadialGradient(
                center: const Alignment(0.0, -0.6),
                radius: 1.2,
                colors: [
                  AegisColors.electricLime.withValues(alpha: 0.03),
                  Colors.transparent,
                ],
              ),
            ),
          ),
        ),
        child,
      ],
    );
  }
}

class _CyberGridPainter extends CustomPainter {
  static const double _gridSize = 28.0;

  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = AegisColors.surfaceBorder.withValues(alpha: 0.35)
      ..strokeWidth = 0.6
      ..style = PaintingStyle.stroke;

    for (double x = 0; x < size.width; x += _gridSize) {
      canvas.drawLine(Offset(x, 0), Offset(x, size.height), paint);
    }
    for (double y = 0; y < size.height; y += _gridSize) {
      canvas.drawLine(Offset(0, y), Offset(size.width, y), paint);
    }
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}
