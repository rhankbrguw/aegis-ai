import 'package:flutter/material.dart';

/// 8-pt Grid Spacing, Radii, and Animation Tokens.
class AegisTokens {
  const AegisTokens._();

  // Spacing (8pt Grid)
  static const double space4 = 4.0;
  static const double space8 = 8.0;
  static const double space12 = 12.0;
  static const double space16 = 16.0;
  static const double space24 = 24.0;
  static const double space32 = 32.0;

  // Corner Radii
  static const double radiusSmall = 8.0;
  static const double radiusMedium = 16.0;
  static const double radiusLarge = 24.0;

  // Animation Durations
  static const Duration animFast = Duration(milliseconds: 150);
  static const Duration animStandard = Duration(milliseconds: 300);
  static const Curve animCurve = Curves.easeOutCubic;
}
