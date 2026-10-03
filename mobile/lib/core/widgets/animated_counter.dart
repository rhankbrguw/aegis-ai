import 'package:flutter/material.dart';

/// Smooth Mechanical Rolling Counter Ticker.
class AnimatedCounter extends StatelessWidget {
  final double value;
  final String prefix;
  final String suffix;
  final int decimalPlaces;
  final TextStyle style;
  final Duration duration;

  const AnimatedCounter({
    super.key,
    required this.value,
    this.prefix = '',
    this.suffix = '',
    this.decimalPlaces = 1,
    required this.style,
    this.duration = const Duration(milliseconds: 600),
  });

  @override
  Widget build(BuildContext context) {
    return TweenAnimationBuilder<double>(
      tween: Tween<double>(begin: 0.0, end: value),
      duration: duration,
      curve: Curves.easeOutCubic,
      builder: (context, animatedValue, child) {
        final formatted = animatedValue.toStringAsFixed(decimalPlaces);
        return Text(
          '$prefix$formatted$suffix',
          style: style,
        );
      },
    );
  }
}
