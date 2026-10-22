import 'package:flutter/material.dart';
import '../constants/colors.dart';
import '../constants/tokens.dart';
import '../constants/typography.dart';

/// Industrial Dark Design System for AegisAI.
class AegisTheme {
  const AegisTheme._();

  static ThemeData get darkTheme {
    return ThemeData(
      brightness: Brightness.dark,
      scaffoldBackgroundColor: AegisColors.background,
      primaryColor: AegisColors.electricLime,
      cardColor: AegisColors.surfacePrimary,
      fontFamily: 'monospace',
      appBarTheme: AppBarTheme(
        backgroundColor: AegisColors.background,
        elevation: 0,
        centerTitle: false,
        titleTextStyle: AegisTypography.heroTitle,
      ),
      cardTheme: CardThemeData(
        color: AegisColors.surfacePrimary,
        elevation: 0,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(AegisTokens.radiusMedium),
          side: const BorderSide(color: AegisColors.surfaceBorder, width: 1),
        ),
      ),
    );
  }
}
