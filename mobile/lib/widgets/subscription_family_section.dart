import 'package:flutter/material.dart';
import '../models/family_member.dart';

class SubscriptionFamilySection extends StatelessWidget {
  final List<FamilyMember> userFamilyMembers;
  final List<FamilyMember> allFamilyMembers;
  final FamilyMember? payerFamilyMember;
  final bool payedByJointAccount;
  final bool isFamilyExpanded;
  final Function(bool) onFamilyExpandedToggle;
  final Function(FamilyMember, bool) onFamilyMemberToggled;
  final Function(bool) onJointAccountToggled;
  final Function(String?) onPayerChanged;
  final bool hasFamilyMembers;

  const SubscriptionFamilySection({
    super.key,
    required this.userFamilyMembers,
    required this.allFamilyMembers,
    required this.payerFamilyMember,
    required this.payedByJointAccount,
    required this.isFamilyExpanded,
    required this.onFamilyExpandedToggle,
    required this.onFamilyMemberToggled,
    required this.onJointAccountToggled,
    required this.onPayerChanged,
    required this.hasFamilyMembers,
  });

  @override
  Widget build(BuildContext context) {
    // If no family members, don't show the section
    if (!hasFamilyMembers) {
      return const SizedBox.shrink();
    }

    return Card(
      elevation: 2,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(16),
      ),
      child: Column(
        children: [
          // Collapsible header
          InkWell(
            onTap: () => onFamilyExpandedToggle(!isFamilyExpanded),
            child: Padding(
              padding: const EdgeInsets.all(16.0),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    'Family',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.w600,
                      color: Theme.of(context).colorScheme.primary,
                    ),
                  ),
                  Icon(
                    isFamilyExpanded
                        ? Icons.keyboard_arrow_up
                        : Icons.keyboard_arrow_down,
                    color: Theme.of(context).colorScheme.primary,
                  ),
                ],
              ),
            ),
          ),

          // Collapsible content
          AnimatedContainer(
            duration: const Duration(milliseconds: 200),
            height: isFamilyExpanded ? null : 0,
            child: isFamilyExpanded
                ? Padding(
                    padding: const EdgeInsets.fromLTRB(
                      16.0,
                      0,
                      16.0,
                      16.0,
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Divider(),
                        const Text(
                          'Link family members to this subscription',
                          style: TextStyle(
                            fontSize: 14,
                            fontWeight: FontWeight.w400,
                          ),
                        ),
                        const SizedBox(height: 16),

                        // Who uses this subscription
                        Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text(
                              'Who uses this subscription',
                              style: TextStyle(
                                fontSize: 16,
                                fontWeight: FontWeight.w500,
                              ),
                            ),
                            const SizedBox(height: 8),
                            Wrap(
                              spacing: 8.0,
                              runSpacing: 8.0,
                              children: _buildFamilyMemberChips(context),
                            ),
                          ],
                        ),

                        const SizedBox(height: 16),

                        // Joint account toggle
                        Row(
                          children: [
                            Switch(
                              value: payedByJointAccount,
                              onChanged: onJointAccountToggled,
                            ),
                            const SizedBox(width: 8),
                            Text(
                              'Paid by joint account',
                              style: TextStyle(
                                fontSize: 16,
                                color: Theme.of(context).textTheme.bodyLarge?.color,
                              ),
                            ),
                          ],
                        ),

                        const SizedBox(height: 16),

                        // Who pays for this subscription (only shown if not paid by joint account)
                        if (!payedByJointAccount)
                          DropdownButtonFormField<String?>(
                            value: payerFamilyMember?.id,
                            decoration: const InputDecoration(
                              labelText: 'Who pays for this subscription',
                              prefixIcon: Icon(Icons.account_balance_wallet),
                            ),
                            items: [
                              const DropdownMenuItem<String?>(
                                value: null,
                                child: Text('None'),
                              ),
                              ...allFamilyMembers
                                  .where((member) => !member.isKid)
                                  .map((member) {
                                    return DropdownMenuItem<String?>(
                                      value: member.id,
                                      child: Text(member.name),
                                    );
                                  }),
                            ],
                            onChanged: onPayerChanged,
                          ),
                      ],
                    ),
                  )
                : const SizedBox.shrink(),
          ),
        ],
      ),
    );
  }

  // Build family member chips for selection
  List<Widget> _buildFamilyMemberChips(BuildContext context) {
    return allFamilyMembers.map((member) {
      final isSelected = userFamilyMembers.any(
        (selectedMember) => selectedMember.id == member.id,
      );
      return FilterChip(
        label: Text(member.name),
        selected: isSelected,
        onSelected: (selected) => onFamilyMemberToggled(member, selected),
        backgroundColor: Theme.of(context).chipTheme.backgroundColor,
        selectedColor: Theme.of(context).colorScheme.primary.withOpacity(0.7),
        checkmarkColor: Colors.white,
      );
    }).toList();
  }
}