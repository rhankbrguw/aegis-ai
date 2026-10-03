import 'package:flutter/material.dart';
import '../../../../core/constants/colors.dart';
import '../../../../core/constants/tokens.dart';
import '../../../../core/constants/typography.dart';
import '../../../../core/widgets/neobrutal_card.dart';

/// Dark Neobrutalist Event Stream Terminal Console with Filter Pills.
class TerminalLogsView extends StatefulWidget {
  final List<String> logs;

  const TerminalLogsView({super.key, required this.logs});

  @override
  State<TerminalLogsView> createState() => _TerminalLogsViewState();
}

class _TerminalLogsViewState extends State<TerminalLogsView> {
  String _selectedFilter = 'ALL';

  @override
  Widget build(BuildContext context) {
    final rawLogs = widget.logs.isEmpty ? _defaultMockLogs : widget.logs;
    final filteredLogs = _filterLogs(rawLogs);

    return NeobrutalCard(
      backgroundColor: AegisColors.terminalBackground,
      padding: const EdgeInsets.all(AegisTokens.space12),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _buildHeader(),
          const SizedBox(height: AegisTokens.space8),
          _buildFilterPills(),
          const Divider(color: AegisColors.surfaceBorder, height: 16),
          _buildLogList(filteredLogs),
        ],
      ),
    );
  }

  Widget _buildHeader() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Row(
          children: [
            _buildDot(AegisColors.laserCrimson),
            const SizedBox(width: 5),
            _buildDot(AegisColors.cyberYellow),
            const SizedBox(width: 5),
            _buildDot(AegisColors.electricLime),
            const SizedBox(width: AegisTokens.space8),
            Text('EVENT LOG STREAM', style: AegisTypography.cardLabel.copyWith(fontSize: 9)),
          ],
        ),
        Text('LIVE WS // 200 OK', style: AegisTypography.actionButton.copyWith(color: AegisColors.electricLime, fontSize: 8)),
      ],
    );
  }

  Widget _buildFilterPills() {
    final filters = ['ALL', 'CACHE', 'FAILOVER'];
    return Row(
      children: filters.map((filter) {
        final isSelected = _selectedFilter == filter;
        return GestureDetector(
          onTap: () => setState(() => _selectedFilter = filter),
          child: Container(
            margin: const EdgeInsets.only(right: 6),
            padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
            decoration: BoxDecoration(
              color: isSelected ? AegisColors.electricCyan : AegisColors.surfaceSecondary,
              borderRadius: BorderRadius.circular(3),
              border: Border.all(color: isSelected ? AegisColors.electricCyan : AegisColors.surfaceBorder, width: 1),
            ),
            child: Text(
              filter,
              style: AegisTypography.microBadge.copyWith(
                color: isSelected ? AegisColors.hardBlack : AegisColors.textMuted,
                fontSize: 8,
              ),
            ),
          ),
        );
      }).toList(),
    );
  }

  Widget _buildLogList(List<String> logs) {
    return SizedBox(
      height: 95,
      child: ListView.builder(
        itemCount: logs.length,
        itemBuilder: (context, index) {
          final log = logs[index];
          final logColor = log.contains('FAILOVER') || log.contains('OUTAGE')
              ? AegisColors.cyberYellow
              : (log.contains('CACHE') ? AegisColors.electricLime : AegisColors.textSecondary);
          return Padding(
            padding: const EdgeInsets.symmetric(vertical: 2.0),
            child: Text(log, style: AegisTypography.terminalLog.copyWith(color: logColor, fontSize: 10, fontWeight: FontWeight.w600)),
          );
        },
      ),
    );
  }

  List<String> _filterLogs(List<String> logs) {
    if (_selectedFilter == 'ALL') return logs;
    return logs.where((l) => l.toUpperCase().contains(_selectedFilter)).toList();
  }

  Widget _buildDot(Color color) {
    return Container(width: 7, height: 7, decoration: BoxDecoration(color: color, shape: BoxShape.circle));
  }

  static const List<String> _defaultMockLogs = [
    r'[11:42:01] POST /v1/chat/completions -> [CACHE_HIT] latency: 1.7ms ($0 cost)',
    r'[11:42:15] POST /v1/chat/completions -> [OPENAI_PRIMARY] latency: 182ms tokens: 45',
    r'[11:42:30] POST /v1/chat/completions -> [CACHE_HIT] latency: 1.5ms ($0 cost)',
    r'[11:43:00] PROBE /health -> [OK] status: healthy redis: connected',
  ];
}
