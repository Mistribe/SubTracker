import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:subscription_tracker/providers/label_provider.dart';
import '../providers/subscription_provider.dart';
import '../providers/family_provider.dart';
import '../widgets/subscription_list.dart';
import 'subscription_form_screen.dart';
import 'settings_screen.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final subscriptionProvider = Provider.of<SubscriptionProvider>(context);
    final labelProvider = Provider.of<LabelProvider>(context);

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
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 12),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Cards row
                  Row(
                    children: [
                      // Monthly cost card
                      Expanded(
                        child: _buildSummaryCard(
                          context,
                          'Monthly',
                          subscriptionProvider.formattedMonthlyCost,
                          Icons.calendar_today,
                          Colors.blue,
                        ),
                      ),
                      const SizedBox(width: 10),

                      // Annual cost card
                      Expanded(
                        child: _buildSummaryCard(
                          context,
                          'Annually',
                          subscriptionProvider.formattedAnnualCost,
                          Icons.calendar_month,
                          Colors.purple,
                        ),
                      ),
                      const SizedBox(width: 10),

                      // Active subscriptions card
                      Expanded(
                        child: _buildSummaryCard(
                          context,
                          'Active',
                          '${subscriptionProvider.activePaymentsCount}${subscriptionProvider.notStartedPaymentsCount > 0 ? ' (${subscriptionProvider.notStartedPaymentsCount})' : ''}',
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
                            '${subscriptionProvider.filteredSubscriptions.length} of ${subscriptionProvider.subscriptions.length}',
                            style: TextStyle(
                              fontSize: 14,
                              color: Theme.of(
                                context,
                              ).colorScheme.onSurface.withValues(alpha: 0.6),
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
                              if (newValue == SubscriptionFilterOption.labels) {
                                _showLabelFilterBottomSheet(
                                  context,
                                  subscriptionProvider,
                                  labelProvider,
                                );
                              } else if (newValue ==
                                  SubscriptionFilterOption.showInactive) {
                                subscriptionProvider.showInactiveSubscriptions =
                                    true;
                              } else if (newValue ==
                                  SubscriptionFilterOption.hideInactive) {
                                subscriptionProvider.showInactiveSubscriptions =
                                    false;
                              } else if (newValue ==
                                  SubscriptionFilterOption.familyMembers) {
                                _showFamilyMemberFilterBottomSheet(
                                  context,
                                  subscriptionProvider,
                                );
                              } else if (newValue ==
                                  SubscriptionFilterOption.payer) {
                                _showPayerFilterBottomSheet(
                                  context,
                                  subscriptionProvider,
                                );
                              }
                            },
                            itemBuilder: (context) => [
                              if (labelProvider.labels.isNotEmpty)
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
                                        'Labels ${subscriptionProvider.selectedLabelIds.isNotEmpty ? "(${subscriptionProvider.selectedLabelIds.length})" : ""}',
                                      ),
                                    ],
                                  ),
                                ),
                              // Family member filter option
                              if (Provider.of<FamilyProvider>(
                                context,
                                listen: false,
                              ).hasFamilyMembers)
                                PopupMenuItem(
                                  value: SubscriptionFilterOption.familyMembers,
                                  child: Row(
                                    children: [
                                      Icon(
                                        Icons.people_outline,
                                        color: Theme.of(
                                          context,
                                        ).colorScheme.primary,
                                      ),
                                      const SizedBox(width: 8),
                                      Text(
                                        'Family Member ${subscriptionProvider.selectedFamilyMemberId != null ? "(1)" : ""}',
                                      ),
                                    ],
                                  ),
                                ),
                              // Payer filter option
                              if (Provider.of<FamilyProvider>(
                                context,
                                listen: false,
                              ).hasFamilyMembers)
                                PopupMenuItem(
                                  value: SubscriptionFilterOption.payer,
                                  child: Row(
                                    children: [
                                      Icon(
                                        Icons.account_balance_wallet_outlined,
                                        color: Theme.of(
                                          context,
                                        ).colorScheme.primary,
                                      ),
                                      const SizedBox(width: 8),
                                      Text(
                                        'Payer ${subscriptionProvider.selectedPayerFamilyMemberId != null ? "(1)" : ""}',
                                      ),
                                    ],
                                  ),
                                ),
                              if (subscriptionProvider
                                  .showInactiveSubscriptions)
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
                              subscriptionProvider.sortOption = newValue;
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
                  const SizedBox(height: 2),
                  // Compact filter bar
                  SingleChildScrollView(
                    scrollDirection: Axis.horizontal,
                    child: Row(children: []),
                  ),
                ],
              ),
            ),

            // Subscription list
            const Expanded(child: SubscriptionList()),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          Navigator.push(
            context,
            MaterialPageRoute(
              builder: (context) => const SubscriptionFormScreen(),
            ),
          );
        },
        tooltip: 'Add Subscription',
        elevation: 4,
        child: const Icon(Icons.add),
      ),
    );
  }

  // Show bottom sheet for family member filtering
  void _showFamilyMemberFilterBottomSheet(
    BuildContext context,
    SubscriptionProvider provider,
  ) {
    final familyMemberProvider = Provider.of<FamilyProvider>(
      context,
      listen: false,
    );

    if (!familyMemberProvider.hasFamilyMembers) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text(
            'No family members available. Add family members in settings.',
          ),
        ),
      );
      return;
    }

    showModalBottomSheet(
      context: context,
      builder: (context) {
        return StatefulBuilder(
          builder: (context, setState) {
            final familyMembers = familyMemberProvider.families;
            final selectedFamilyMemberId = provider.selectedFamilyMemberId;

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
                          'Filter by Family Member',
                          style: Theme.of(context).textTheme.titleLarge,
                        ),
                        if (selectedFamilyMemberId != null)
                          TextButton(
                            onPressed: () {
                              provider.clearFamilyMemberFilter();
                              setState(() {});
                            },
                            child: const Text('Clear'),
                          ),
                      ],
                    ),

                    const Divider(),

                    // Family member list
                    ListView.builder(
                      shrinkWrap: true,
                      itemCount: familyMembers.length,
                      itemBuilder: (context, index) {
                        final familyMember = familyMembers[index];
                        final isSelected =
                            selectedFamilyMemberId == familyMember.id;

                        return RadioListTile<String>(
                          title: Text(familyMember.name),
                          value: familyMember.id,
                          groupValue: selectedFamilyMemberId,
                          onChanged: (value) {
                            provider.selectedFamilyMemberId = value;
                            setState(() {});
                          },
                          selected: isSelected,
                        );
                      },
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

  // Show bottom sheet for payer filtering
  void _showPayerFilterBottomSheet(
    BuildContext context,
    SubscriptionProvider provider,
  ) {
    final familyMemberProvider = Provider.of<FamilyProvider>(
      context,
      listen: false,
    );

    if (!familyMemberProvider.hasFamilyMembers) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text(
            'No family members available. Add family members in settings.',
          ),
        ),
      );
      return;
    }

    showModalBottomSheet(
      context: context,
      builder: (context) {
        return StatefulBuilder(
          builder: (context, setState) {
            final familyMembers = familyMemberProvider.families;
            final selectedPayerFamilyMemberId =
                provider.selectedPayerFamilyMemberId;

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
                          'Filter by Payer',
                          style: Theme.of(context).textTheme.titleLarge,
                        ),
                        if (selectedPayerFamilyMemberId != null)
                          TextButton(
                            onPressed: () {
                              provider.clearPayerFamilyMemberFilter();
                              setState(() {});
                            },
                            child: const Text('Clear'),
                          ),
                      ],
                    ),

                    const Divider(),

                    // Payer list with Family (common account) option
                    Column(
                      children: [
                        // Family (common account) option
                        RadioListTile<String>(
                          title: const Text('Family (common account)'),
                          subtitle: const Text(
                            'Subscriptions with no specific payer',
                          ),
                          value: kFamilyCommonAccountId,
                          groupValue: selectedPayerFamilyMemberId,
                          onChanged: (value) {
                            provider.selectedPayerFamilyMemberId = value;
                            setState(() {});
                          },
                          selected:
                              selectedPayerFamilyMemberId ==
                              kFamilyCommonAccountId,
                        ),

                        const Divider(),

                        // Individual family members
                        ListView.builder(
                          shrinkWrap: true,
                          physics: const NeverScrollableScrollPhysics(),
                          itemCount: familyMembers.length,
                          itemBuilder: (context, index) {
                            final familyMember = familyMembers[index];
                            final isSelected =
                                selectedPayerFamilyMemberId == familyMember.id;

                            return RadioListTile<String>(
                              title: Text(familyMember.name),
                              value: familyMember.id,
                              groupValue: selectedPayerFamilyMemberId,
                              onChanged: (value) {
                                provider.selectedPayerFamilyMemberId = value;
                                setState(() {});
                              },
                              selected: isSelected,
                            );
                          },
                        ),
                      ],
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

  // Show bottom sheet for label filtering - simplified implementation
  void _showLabelFilterBottomSheet(
    BuildContext context,
    SubscriptionProvider subscriptionProvider,
    LabelProvider labelProvider,
  ) {
    showModalBottomSheet(
      context: context,
      builder: (context) {
        return StatefulBuilder(
          builder: (context, setState) {
            final allLabels = labelProvider.labels;
            final selectedLabelIds = subscriptionProvider.selectedLabelIds;

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
                              subscriptionProvider.clearLabelFilters();
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
                              subscriptionProvider.toggleLabelFilter(label.id);
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
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
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
                  color: iconColor.withValues(alpha: isDark ? 0.8 : 1.0),
                ),
                const SizedBox(width: 6),
                Text(
                  title,
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w500,
                    color: Theme.of(
                      context,
                    ).colorScheme.onSurface.withValues(alpha: 0.7),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 6),
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
