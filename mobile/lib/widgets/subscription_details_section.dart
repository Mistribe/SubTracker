import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class SubscriptionDetailsSection extends StatelessWidget {
  final String selectedDuration;
  final TextEditingController customMonthsController;
  final TextEditingController recurrenceCountController;
  final TextEditingController freeTrialMonthsController;
  final DateTime startDate;
  final bool useRecurrenceCount;
  final bool useFreeTrial;
  final bool isEditMode;
  final bool isActiveSubscription;
  final Function(String?) onDurationChanged;
  final Function(String) onCustomMonthsChanged;
  final Function(BuildContext) onSelectStartDate;
  final Function(bool) onRecurrenceCountToggled;
  final Function(bool) onFreeTrialToggled;
  final Function(String) onFreeTrialMonthsChanged;
  final String durationLabel;

  const SubscriptionDetailsSection({
    Key? key,
    required this.selectedDuration,
    required this.customMonthsController,
    required this.recurrenceCountController,
    required this.freeTrialMonthsController,
    required this.startDate,
    required this.useRecurrenceCount,
    required this.useFreeTrial,
    required this.isEditMode,
    required this.isActiveSubscription,
    required this.onDurationChanged,
    required this.onCustomMonthsChanged,
    required this.onSelectStartDate,
    required this.onRecurrenceCountToggled,
    required this.onFreeTrialToggled,
    required this.onFreeTrialMonthsChanged,
    required this.durationLabel,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    // Only show if not in edit mode or if in edit mode and subscription is active
    if (isEditMode && !isActiveSubscription) {
      return const SizedBox.shrink();
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Padding(
          padding: const EdgeInsets.only(left: 4, bottom: 12),
          child: Text(
            'Subscription Details',
            style: TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.w600,
              color: Theme.of(context).colorScheme.primary,
            ),
          ),
        ),
        Card(
          elevation: 2,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(16),
          ),
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              children: [
                // Recurrence type dropdown
                DropdownButtonFormField<String>(
                  value: selectedDuration,
                  decoration: const InputDecoration(
                    labelText: 'Recurrence Type',
                    prefixIcon: Icon(Icons.repeat),
                  ),
                  items: const [
                    DropdownMenuItem(
                      value: '1 month',
                      child: Text('Monthly'),
                    ),
                    DropdownMenuItem(
                      value: '3 months',
                      child: Text('Quarterly (3 months)'),
                    ),
                    DropdownMenuItem(
                      value: '6 months',
                      child: Text('Semi-Annual (6 months)'),
                    ),
                    DropdownMenuItem(
                      value: '12 months',
                      child: Text('Annual (12 months)'),
                    ),
                    DropdownMenuItem(
                      value: 'Custom',
                      child: Text('Custom'),
                    ),
                  ],
                  onChanged: onDurationChanged,
                ),

                // Custom months input (shown only if custom is selected)
                if (selectedDuration == 'Custom')
                  Padding(
                    padding: const EdgeInsets.only(top: 16.0),
                    child: TextFormField(
                      controller: customMonthsController,
                      decoration: const InputDecoration(
                        labelText: 'Custom Months',
                        hintText: 'Enter number of months',
                        prefixIcon: Icon(Icons.calendar_today),
                      ),
                      keyboardType: TextInputType.number,
                      inputFormatters: [
                        FilteringTextInputFormatter.digitsOnly,
                      ],
                      validator: (value) {
                        if (selectedDuration == 'Custom') {
                          if (value == null || value.isEmpty) {
                            return 'Please enter number of months';
                          }
                          final months = int.tryParse(value);
                          if (months == null || months <= 0) {
                            return 'Please enter a valid number';
                          }
                        }
                        return null;
                      },
                      onChanged: onCustomMonthsChanged,
                    ),
                  ),

                const SizedBox(height: 16),

                // Start date picker
                Container(
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(
                      color: Theme.of(context).dividerColor.withOpacity(0.5),
                    ),
                  ),
                  child: ListTile(
                    leading: Icon(
                      Icons.calendar_today,
                      color: Colors.blue.withOpacity(
                        Theme.of(context).brightness == Brightness.dark ? 0.8 : 1.0,
                      ),
                    ),
                    title: const Text('Start Date'),
                    subtitle: Text(
                      '${startDate.month}/${startDate.day}/${startDate.year}',
                    ),
                    onTap: () => onSelectStartDate(context),
                  ),
                ),

                const SizedBox(height: 16),

                // Free trial toggle
                Row(
                  children: [
                    Switch(
                      value: useFreeTrial,
                      onChanged: onFreeTrialToggled,
                    ),
                    const SizedBox(width: 8),
                    Text(
                      'Include free trial period',
                      style: TextStyle(
                        fontSize: 16,
                        color: Theme.of(context).textTheme.bodyLarge?.color,
                      ),
                    ),
                  ],
                ),

                // Free trial months input (shown only if toggle is on)
                if (useFreeTrial)
                  Padding(
                    padding: const EdgeInsets.only(
                      bottom: 16.0,
                      top: 8.0,
                    ),
                    child: TextFormField(
                      controller: freeTrialMonthsController,
                      decoration: const InputDecoration(
                        labelText: 'Free Trial Months',
                        hintText: 'Enter number of months',
                        prefixIcon: Icon(Icons.card_giftcard),
                      ),
                      keyboardType: TextInputType.number,
                      inputFormatters: [
                        FilteringTextInputFormatter.digitsOnly,
                      ],
                      validator: (value) {
                        if (useFreeTrial) {
                          if (value == null || value.isEmpty) {
                            return 'Please enter number of months';
                          }
                          final months = int.tryParse(value);
                          if (months == null || months <= 0) {
                            return 'Please enter a valid number';
                          }
                        }
                        return null;
                      },
                      onChanged: onFreeTrialMonthsChanged,
                    ),
                  ),

                const SizedBox(height: 16),

                // Recurrence count toggle
                Row(
                  children: [
                    Switch(
                      value: useRecurrenceCount,
                      onChanged: onRecurrenceCountToggled,
                    ),
                    const SizedBox(width: 8),
                    Text(
                      'Specify number of recurrences',
                      style: TextStyle(
                        fontSize: 16,
                        color: Theme.of(context).textTheme.bodyLarge?.color,
                      ),
                    ),
                  ],
                ),

                // Recurrence count input (shown only if toggle is on)
                if (useRecurrenceCount)
                  Padding(
                    padding: const EdgeInsets.only(
                      bottom: 8.0,
                      top: 8.0,
                    ),
                    child: TextFormField(
                      controller: recurrenceCountController,
                      decoration: InputDecoration(
                        labelText: 'Number of $durationLabel',
                        hintText: 'Enter number of recurrences',
                        prefixIcon: const Icon(Icons.repeat),
                      ),
                      keyboardType: TextInputType.number,
                      inputFormatters: [
                        FilteringTextInputFormatter.digitsOnly,
                      ],
                      validator: (value) {
                        if (useRecurrenceCount) {
                          if (value == null || value.isEmpty) {
                            return 'Please enter number of recurrences';
                          }
                          final count = int.tryParse(value);
                          if (count == null || count <= 0) {
                            return 'Please enter a valid number';
                          }
                        }
                        return null;
                      },
                    ),
                  ),
              ],
            ),
          ),
        ),
      ],
    );
  }
}
