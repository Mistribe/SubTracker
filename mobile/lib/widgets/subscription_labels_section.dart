import 'package:flutter/material.dart';
import '../models/label.dart';

class SubscriptionLabelsSection extends StatelessWidget {
  final List<Label> selectedLabels;
  final List<Label> allLabels;
  final bool isLabelsExpanded;
  final Function(bool) onLabelsExpandedToggle;
  final Function(Label, bool) onLabelToggled;
  final VoidCallback onAddCustomLabel;

  const SubscriptionLabelsSection({
    super.key,
    required this.selectedLabels,
    required this.allLabels,
    required this.isLabelsExpanded,
    required this.onLabelsExpandedToggle,
    required this.onLabelToggled,
    required this.onAddCustomLabel,
  });

  @override
  Widget build(BuildContext context) {
    return Card(
      elevation: 2,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(16),
      ),
      child: Column(
        children: [
          // Collapsible header
          InkWell(
            onTap: () => onLabelsExpandedToggle(!isLabelsExpanded),
            child: Padding(
              padding: const EdgeInsets.all(16.0),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    'Labels',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.w600,
                      color: Theme.of(context).colorScheme.primary,
                    ),
                  ),
                  Icon(
                    isLabelsExpanded
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
            height: isLabelsExpanded ? null : 0,
            child: isLabelsExpanded
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
                          'Categorize your subscription',
                          style: TextStyle(
                            fontSize: 14,
                            fontWeight: FontWeight.w400,
                          ),
                        ),
                        const SizedBox(height: 16),
                        Wrap(
                          spacing: 8.0,
                          runSpacing: 8.0,
                          children: _buildLabelChips(context),
                        ),
                        Align(
                          alignment: Alignment.centerRight,
                          child: TextButton.icon(
                            icon: const Icon(Icons.add),
                            label: const Text('Add Custom Label'),
                            onPressed: onAddCustomLabel,
                          ),
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

  // Build label chips for selection
  List<Widget> _buildLabelChips(BuildContext context) {
    return allLabels.map((label) {
      final isSelected = selectedLabels.any(
        (selectedLabel) => selectedLabel.id == label.id,
      );
      return FilterChip(
        label: Text(label.name),
        selected: isSelected,
        onSelected: (selected) => onLabelToggled(label, selected),
        backgroundColor: Color(
          int.parse(label.color.substring(1, 7), radix: 16) + 0xFF000000,
        ).withOpacity(0.2),
        selectedColor: Color(
          int.parse(label.color.substring(1, 7), radix: 16) + 0xFF000000,
        ).withOpacity(0.7),
        checkmarkColor: Colors.white,
      );
    }).toList();
  }
}