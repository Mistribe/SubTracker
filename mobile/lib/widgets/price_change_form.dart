import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';

import '../models/subscription.dart';
import '../models/currency.dart';
import '../providers/subscription_provider.dart';

class PriceChangeForm extends StatefulWidget {
  final Subscription subscription;

  const PriceChangeForm({super.key, required this.subscription});

  @override
  State<PriceChangeForm> createState() => _PriceChangeFormState();
}

class _PriceChangeFormState extends State<PriceChangeForm> {
  final _formKey = GlobalKey<FormState>();
  late final _priceController;
  DateTime _selectedDate = DateTime.now();
  late String _selectedCurrency;

  // Get currencies from the Currency enum
  final List<String> _currencies = Currency.codes;

  @override
  void dispose() {
    _priceController.dispose();
    super.dispose();
  }

  @override
  void initState() {
    super.initState();

    final currentDetail = widget.subscription.getLastPaymentDetail();
    _priceController = TextEditingController(
      text: currentDetail.price.toString(),
    );
    _selectedCurrency = currentDetail.currency;
  }

  Future<void> _selectEffectiveDate(BuildContext context) async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate: _selectedDate,
      firstDate: widget.subscription.getLastPaymentDetail().startDate,
      lastDate: DateTime(2101),
    );
    if (picked != null && picked != _selectedDate) {
      setState(() {
        _selectedDate = picked;
      });
    }
  }

  Future<void> _submitForm() async {
    if (_formKey.currentState!.validate()) {
      final price = double.tryParse(_priceController.text) ?? 0.0;

      try {
        // Show loading indicator
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Updating subscription...'),
            duration: Duration(seconds: 1),
          ),
        );

        // Update the subscription detail
        final currentDetail = widget.subscription.getLastPaymentDetail();
        await Provider.of<SubscriptionProvider>(
          context,
          listen: false,
        ).addPaymentDetailEntry(
          widget.subscription.id,
          price,
          _selectedDate,
          endDate: null,
          months: currentDetail.months,
        );

        // Show success message
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Price updated successfully'),
            duration: Duration(seconds: 2),
          ),
        );

        // Close the form
        Navigator.of(context).pop();
      } catch (e) {
        // Show error message
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Error updating subscription: ${e.toString()}'),
            backgroundColor: Colors.red,
            duration: const Duration(seconds: 3),
          ),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    // Get the keyboard height to adjust the bottom padding
    final mediaQuery = MediaQuery.of(context);
    final keyboardHeight = mediaQuery.viewInsets.bottom;

    return Padding(
      padding: EdgeInsets.only(
        top: 16,
        left: 16,
        right: 16,
        bottom: keyboardHeight + 16,
      ),
      child: SingleChildScrollView(
        child: Form(
          key: _formKey,
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              const Text(
                'Add a change of price',
                style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 16),
              Row(
                children: [
                  Expanded(
                    flex: 7,
                    child: TextFormField(
                      controller: _priceController,
                      decoration: InputDecoration(
                        labelText: 'New Price',
                        hintText: 'Enter the new price',
                        prefixIcon: Icon(_selectedCurrency == 'USD' ? Icons.attach_money : Icons.currency_exchange),
                      ),
                      keyboardType: const TextInputType.numberWithOptions(
                        decimal: true,
                      ),
                      inputFormatters: [
                        FilteringTextInputFormatter.allow(RegExp(r'^\d+\.?\d{0,2}')),
                      ],
                    ),
                  ),
                  const SizedBox(width: 8),
                  Expanded(
                    flex: 3,
                    child: DropdownButtonFormField<String>(
                      value: _selectedCurrency,
                      decoration: const InputDecoration(
                        labelText: 'Currency',
                        contentPadding: EdgeInsets.symmetric(horizontal: 10, vertical: 15),
                      ),
                      items: _currencies.map<DropdownMenuItem<String>>((String value) {
                        return DropdownMenuItem<String>(
                          value: value,
                          child: Text(value),
                        );
                      }).toList(),
                      onChanged: (String? newValue) {
                        if (newValue != null) {
                          setState(() {
                            _selectedCurrency = newValue;
                          });
                        }
                      },
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              ListTile(
                leading: const Icon(Icons.calendar_today),
                title: const Text('Effective Date'),
                subtitle: Text(
                  '${_selectedDate.month}/${_selectedDate.day}/${_selectedDate.year}',
                ),
                onTap: () => _selectEffectiveDate(context),
              ),
              const SizedBox(height: 24),
              ElevatedButton(
                onPressed: _submitForm,
                style: ElevatedButton.styleFrom(
                  padding: const EdgeInsets.symmetric(vertical: 12),
                ),
                child: const Text(
                  'Update Payment',
                  style: TextStyle(fontSize: 16),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
