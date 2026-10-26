import 'package:device_preview/device_preview.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'core/constants/strings.dart';
import 'core/theme/theme.dart';
import 'features/dashboard/presentation/dashboard_view.dart';

void main() {
  WidgetsFlutterBinding.ensureInitialized();
  runApp(
    ProviderScope(
      child: DevicePreview(
        enabled: !kReleaseMode,
        builder: (context) => const AegisApp(),
      ),
    ),
  );
}

class AegisApp extends StatelessWidget {
  const AegisApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      locale: DevicePreview.locale(context),
      builder: DevicePreview.appBuilder,
      title: AegisStrings.appTitle,
      debugShowCheckedModeBanner: false,
      theme: AegisTheme.darkTheme,
      home: const DashboardView(),
    );
  }
}
