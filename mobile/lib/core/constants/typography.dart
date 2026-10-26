import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'colors.dart';

/// Modern Minimalist & Neobrutalist Typography Tokens.
class AegisTypography {
  const AegisTypography._();

  /// Hero Title (Space Grotesk - Modern Tech / Brutal Aesthetic)
  static TextStyle get heroTitle => GoogleFonts.spaceGrotesk(
        color: AegisColors.textPrimary,
        fontSize: 16,
        fontWeight: FontWeight.w800,
        letterSpacing: 2.2,
      );

  /// Section & Header Titles
  static TextStyle get sectionTitle => GoogleFonts.spaceGrotesk(
        color: AegisColors.textPrimary,
        fontSize: 14,
        fontWeight: FontWeight.w700,
        letterSpacing: 1.2,
      );

  /// Micro-Labels & Badges (Space Grotesk Uppercase)
  static TextStyle get microBadge => GoogleFonts.spaceGrotesk(
        color: AegisColors.hardBlack,
        fontSize: 10,
        fontWeight: FontWeight.w900,
        letterSpacing: 1.2,
      );

  /// Card Headers & Labels
  static TextStyle get cardLabel => GoogleFonts.spaceGrotesk(
        color: AegisColors.textMuted,
        fontSize: 10,
        fontWeight: FontWeight.w900,
        letterSpacing: 1.4,
      );

  /// Mechanical & FinOps Telemetry Metrics (JetBrains Mono)
  static TextStyle get telemetryMetric => GoogleFonts.jetBrainsMono(
        color: AegisColors.electricLime,
        fontSize: 20,
        fontWeight: FontWeight.w800,
        letterSpacing: -0.5,
      );

  /// Realtime Terminal Logs
  static TextStyle get terminalLog => GoogleFonts.jetBrainsMono(
        color: AegisColors.textSecondary,
        fontSize: 11,
        height: 1.45,
        fontWeight: FontWeight.w500,
      );

  /// Tactical Action Text
  static TextStyle get actionButton => GoogleFonts.spaceGrotesk(
        color: AegisColors.textPrimary,
        fontSize: 11,
        fontWeight: FontWeight.w900,
        letterSpacing: 1.1,
      );
}
