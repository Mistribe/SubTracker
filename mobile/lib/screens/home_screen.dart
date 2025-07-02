import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/subscription_provider.dart';
import '../providers/theme_provider.dart';
import '../widgets/subscription_list.dart';
import 'subscription_form_screen.dart';
import 'settings_screen.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final paymentProvider = Provider.of<SubscriptionProvider>(context);
    final themeProvider = Provider.of<ThemeProvider>(context);
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Scaffold(
      appBar: AppBar(
        title: const Text(
          'Subscription Tracker',
          style: TextStyle(fontSize: 20, fontWeight: FontWeight.w600),
        ),
        elevation: 0,
        actions: [
          IconButton(
            icon: const Icon(Icons.settings),
            onPressed: () {
              Navigator.push(
                context,
                MaterialPageRoute(builder: (context) => const SettingsScreen()),
              );
            },
            tooltip: 'Settings',
          ),
        ],
      ),
      body: SafeArea(
        child: Column(
          children: [
            // Summary cards with modern design
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Title for the summary section
                  Padding(
                    padding: const EdgeInsets.only(left: 4, bottom: 12),
                    child: Text(
                      'Summary',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.w600,
                        color: Theme.of(context).colorScheme.primary,
                      ),
                    ),
                  ),

                  // Cards row
                  Row(
                    children: [
                      // Monthly cost card
                      Expanded(
                        child: _buildSummaryCard(
                          context,
                          'Monthly',
                          paymentProvider.formattedMonthlyCost,
                          Icons.calendar_today,
                          Colors.blue,
                        ),
                      ),
                      const SizedBox(width: 12),

                      // Annual cost card
                      Expanded(
                        child: _buildSummaryCard(
                          context,
                          'Annually',
                          paymentProvider.formattedAnnualCost,
                          Icons.calendar_month,
                          Colors.purple,
                        ),
                      ),
                      const SizedBox(width: 12),

                      // Active subscriptions card
                      Expanded(
                        child: _buildSummaryCard(
                          context,
                          'Active',
                          '${paymentProvider.activePaymentsCount}${paymentProvider.notStartedPaymentsCount > 0 ? ' (${paymentProvider.notStartedPaymentsCount})' : ''}',
                          Icons.check_circle,
                          Colors.green,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),

            // Subscriptions section title
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20),
              child: Column(
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        'Your Subscriptions',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.w600,
                          color: Theme.of(context).colorScheme.primary,
                        ),
                      ),
                      Row(
                        children: [
                          Text(
                            '${paymentProvider.filteredSubscriptions.length} of ${paymentProvider.subscriptions.length}',
                            style: TextStyle(
                              fontSize: 14,
                              color: Theme.of(
                                context,
                              ).colorScheme.onSurface.withOpacity(0.6),
                            ),
                          ),
                          // Filter buttons
                          PopupMenuButton<SubscriptionFilterOption>(
                            tooltip: 'Filter',
                            icon: Icon(
                              Icons.filter_alt,
                              color: Theme.of(context).colorScheme.primary,
                              size: 20,
                            ),
                            onSelected: (SubscriptionFilterOption newValue) {
                              switch (newValue) {
                                case SubscriptionFilterOption.labels:
                                  _showLabelFilterBottomSheet(
                                    context,
                                    paymentProvider,
                                  );
                                  break;
                                case SubscriptionFilterOption.showInactive:
                                  paymentProvider.showInactiveSubscriptions =
                                      true;
                                  break;
                                case SubscriptionFilterOption.hideInactive:
                                  paymentProvider.showInactiveSubscriptions =
                                      false;
                                  break;
                              }
                            },
                            itemBuilder: (context) => [
                              if (paymentProvider.labels.isNotEmpty)
                                PopupMenuItem(
                                  value: SubscriptionFilterOption.labels,
                                  child: Row(
                                    children: [
                                      Icon(
                                        Icons.label_outline,
                                        color: Theme.of(
                                          context,
                                        ).colorScheme.primary,
                                      ),
                                      const SizedBox(width: 8),
                                      Text(
                                        'Labels ${paymentProvider.selectedLabelIds.isNotEmpty ? "(${paymentProvider.selectedLabelIds.length})" : ""}',
                                      ),
                                    ],
                                  ),
                                ),
                              if (paymentProvider.showInactiveSubscriptions)
                                PopupMenuItem(
                                  value: SubscriptionFilterOption.hideInactive,
                                  child: Row(
                                    children: [
                                      Icon(
                                        Icons.play_disabled,
                                        color: Theme.of(
                                          context,
                                        ).colorScheme.primary,
                                      ),
                                      const SizedBox(width: 8),
                                      const Text('Hide Inactive'),
                                    ],
                                  ),
                                )
                              else
                                PopupMenuItem(
                                  value: SubscriptionFilterOption.showInactive,
                                  child: Row(
                                    children: [
                                      Icon(
                                        Icons.play_arrow,
                                        color: Theme.of(
                                          context,
                                        ).colorScheme.primary,
                                      ),
                                      const SizedBox(width: 8),
                                      const Text('Show Inactive'),
                                    ],
                                  ),
                                ),
                            ],
                          ),
                          // Sort buttons
                          PopupMenuButton<SubscriptionSortOption>(
                            tooltip: 'Sort',
                            icon: Icon(
                              Icons.sort,
                              color: Theme.of(context).colorScheme.primary,
                              size: 20,
                            ),
                            onSelected: (SubscriptionSortOption newValue) {
                              paymentProvider.sortOption = newValue;
                            },
                            itemBuilder: (context) => [
                              const PopupMenuItem(
                                value: SubscriptionSortOption.none,
                                child: Text('Default order'),
                              ),
                              const PopupMenuItem(
                                value: SubscriptionSortOption.nameAsc,
                                child: Text('Name (A-Z)'),
                              ),
                              const PopupMenuItem(
                                value: SubscriptionSortOption.nameDesc,
                                child: Text('Name (Z-A)'),
                              ),
                              const PopupMenuItem(
                                value: SubscriptionSortOption.nextPaymentAsc,
                                child: Text('Next payment (earliest)'),
                              ),
                              const PopupMenuItem(
                                value: SubscriptionSortOption.nextPaymentDesc,
                                child: Text('Next payment (latest)'),
                              ),
                            ],
                          ),
                        ],
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),
                  // Compact filter bar
                  SingleChildScrollView(
                    scrollDirection: Axis.horizontal,
                    child: Row(children: []),
                  ),
                ],
              ),
            ),

            const SizedBox(height: 8),

            // Subscription list
            const Expanded(child: SubscriptionList()),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () {
          Navigator.push(
            context,
            MaterialPageRoute(
              builder: (context) => const SubscriptionFormScreen(),
            ),
          );
        },
        tooltip: 'Add Subscription',
        icon: const Icon(Icons.add),
        label: const Text('Add', style: TextStyle(fontSize: 16)),
        elevation: 4,
      ),
    );
  }

  // Show bottom sheet for label filtering - simplified implementation
  void _showLabelFilterBottomSheet(
    BuildContext context,
    SubscriptionProvider provider,
  ) {
    showModalBottomSheet(
      context: context,
      builder: (context) {
        return StatefulBuilder(
          builder: (context, setState) {
            final allLabels = provider.labels;
            final selectedLabelIds = provider.selectedLabelIds;

            return SafeArea(
              child: Padding(
                padding: const EdgeInsets.symmetric(
                  horizontal: 16,
                  vertical: 12,
                ),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Header with clear button
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text(
                          'Filter by Label',
                          style: Theme.of(context).textTheme.titleLarge,
                        ),
                        if (selectedLabelIds.isNotEmpty)
                          TextButton(
                            onPressed: () {
                              provider.clearLabelFilters();
                              setState(() {});
                            },
                            child: const Text('Clear'),
                          ),
                      ],
                    ),

                    const Divider(),

                    // Label chips
                    if (allLabels.isEmpty)
                      const Center(
                        child: Padding(
                          padding: EdgeInsets.all(16),
                          child: Text('No labels available'),
                        ),
                      )
                    else
                      Wrap(
                        spacing: 8,
                        runSpacing: 8,
                        children: allLabels.map((label) {
                          final isSelected = selectedLabelIds.contains(
                            label.id,
                          );
                          final labelColor = Color(
                            int.parse(label.color.substring(1, 7), radix: 16) +
                                0xFF000000,
                          );

                          return FilterChip(
                            selected: isSelected,
                            label: Text(label.name),
                            onSelected: (_) {
                              provider.toggleLabelFilter(label.id);
                              setState(() {});
                            },
                            selectedColor: labelColor,
                            checkmarkColor: Colors.white,
                          );
                        }).toList(),
                      ),

                    const SizedBox(height: 16),

                    // Done button
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton(
                        onPressed: () => Navigator.pop(context),
                        child: const Text('Done'),
                      ),
                    ),
                  ],
                ),
              ),
            );
          },
        );
      },
    );
  }

  Widget _buildSummaryCard(
    BuildContext context,
    String title,
    String value,
    IconData icon,
    Color iconColor,
  ) {
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Card(
      elevation: 2,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 16, horizontal: 12),
        child: Column(
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(
                  icon,
                  size: 16,
                  color: iconColor.withOpacity(isDark ? 0.8 : 1.0),
                ),
                const SizedBox(width: 6),
                Text(
                  title,
                  style: TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w500,
                    color: Theme.of(
                      context,
                    ).colorScheme.onSurface.withOpacity(0.7),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              value,
              style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }
}
