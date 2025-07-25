import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import '../models/currency.dart';

class SubscriptionBasicInfoSection extends StatelessWidget {
  final TextEditingController nameController;
  final TextEditingController priceController;
  final String selectedCurrency;
  final List<String> currencies;
  final bool isEditMode;
  final bool isActiveSubscription;
  final Function(String?) onCurrencyChanged;

  const SubscriptionBasicInfoSection({
    super.key,
    required this.nameController,
    required this.priceController,
    required this.selectedCurrency,
    required this.currencies,
    required this.isEditMode,
    required this.isActiveSubscription,
    required this.onCurrencyChanged,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Padding(
          padding: const EdgeInsets.only(left: 4, bottom: 12),
          child: Text(
            'Basic Information',
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
                TextFormField(
                  controller: nameController,
                  decoration: const InputDecoration(
                    labelText: 'Subscription Name',
                    hintText: 'e.g., Netflix, Gym Membership',
                    prefixIcon: Icon(Icons.payment),
                  ),
                  textCapitalization: TextCapitalization.words,
                  validator: (value) {
                    if (value == null || value.trim().isEmpty) {
                      return 'Please enter a name';
                    }
                    return null;
                  },
                ),
                if (!isEditMode || (isEditMode && isActiveSubscription))
                  const SizedBox(height: 16),
                if (!isEditMode || (isEditMode && isActiveSubscription))
                  Row(
                    children: [
                      Expanded(
                        flex: 7,
                        child: TextFormField(
                          controller: priceController,
                          decoration: InputDecoration(
                            labelText: 'Price',
                            hintText: 'Enter the amount',
                            prefixIcon: Icon(
                              selectedCurrency == Currency.USD.code
                                  ? Icons.attach_money
                                  : Icons.currency_exchange,
                            ),
                          ),
                          enabled:
                              !isEditMode ||
                              (isEditMode && isActiveSubscription),
                          keyboardType: const TextInputType.numberWithOptions(
                            decimal: true,
                          ),
                          inputFormatters: [
                            FilteringTextInputFormatter.allow(
                              RegExp(r'^\d+\.?\d{0,2}'),
                            ),
                          ],
                          validator: (value) {
                            if (value == null || value.isEmpty) {
                              return 'Please enter a price';
                            }
                            try {
                              final price = double.parse(value);
                              if (price <= 0) {
                                return 'Price must be greater than zero';
                              }
                            } catch (e) {
                              return 'Please enter a valid number';
                            }
                            return null;
                          },
                        ),
                      ),
                      const SizedBox(width: 8),
                      Expanded(
                        flex: 3,
                        child: DropdownButtonFormField<String>(
                          value: selectedCurrency,
                          decoration: const InputDecoration(
                            labelText: 'Currency',
                            contentPadding: EdgeInsets.symmetric(
                              horizontal: 10,
                              vertical: 15,
                            ),
                          ),
                          items: currencies.map<DropdownMenuItem<String>>((
                            String value,
                          ) {
                            return DropdownMenuItem<String>(
                              value: value,
                              child: Text(value),
                            );
                          }).toList(),
                          onChanged:
                              !isEditMode ||
                                  (isEditMode && isActiveSubscription)
                              ? onCurrencyChanged
                              : null,
                        ),
                      ),
                    ],
                  ),
              ],
            ),
          ),
        ),
      ],
    );
  }
}
