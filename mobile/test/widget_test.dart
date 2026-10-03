import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:aegis_mobile/main.dart';

void main() {
  testWidgets('Renders AegisAI Tactical SRE Dashboard with Motion Widgets', (WidgetTester tester) async {
    await tester.pumpWidget(
      const ProviderScope(
        child: AegisApp(),
      ),
    );

    await tester.pump(const Duration(milliseconds: 100));

    expect(find.text('AEGIS // SENTINEL CORE'), findsOneWidget);
    expect(find.text('api.rhankbrguw.xyz'), findsOneWidget);
    expect(find.text('BURN VELOCITY'), findsOneWidget);
    expect(find.text('FINOPS SAVINGS'), findsOneWidget);
    expect(find.text('SLIDE TO EMERGENCY HALT >>>'), findsOneWidget);
  });
}
