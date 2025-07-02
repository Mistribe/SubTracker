import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import '../providers/subscription_provider.dart';
import '../models/label.dart';
import '../models/currency.dart';

class AddSubscriptionScreen extends StatefulWidget {
  const AddSubscriptionScreen({super.key});

  @override
  State<AddSubscriptionScreen> createState() => _AddSubscriptionScreenState();
}

class _AddSubscriptionScreenState extends State<AddSubscriptionScreen> {
  final _formKey = GlobalKey<FormState>();
  late int _months = 1;
  final _nameController = TextEditingController();
  final _priceController = TextEditingController();
  final _customMonthsController = TextEditingController();
  String _selectedDuration = '1 month'; // Default value
  DateTime _selectedStartDate = DateTime.now();
  DateTime? _selectedEndDate;
  String _selectedCurrency = Currency.USD.code; // Default currency
  final List<Label> _selectedLabels = [];

  // Get currencies from the Currency enum
  final List<String> _currencies = Currency.codes;

  @override
  void initState() {
    super.initState();
    // Initialize with the default currency from the provider
    WidgetsBinding.instance.addPostFrameCallback((_) {
      final provider = Provider.of<SubscriptionProvider>(context, listen: false);
      setState(() {
        _selectedCurrency = provider.defaultCurrency;
      });
    });
  }

  @override
  void dispose() {
    _nameController.dispose();
    _priceController.dispose();
    _customMonthsController.dispose();
    super.dispose();
  }

  // Show date picker
  Future<void> _selectStartDate(BuildContext context) async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate: _selectedStartDate,
      firstDate: DateTime(2000),
      lastDate: DateTime(2101),
    );
    if (picked != null && picked != _selectedStartDate) {
      setState(() {
        _selectedStartDate = picked;
      });
    }
  }

  // Date picker for end date
  Future<void> _selectEndDate(BuildContext context) async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate: _selectedEndDate ?? DateTime.now(),
      firstDate: DateTime(1900),
      lastDate: DateTime(2101),
    );
    if (picked != null && picked != _selectedEndDate) {
      setState(() {
        if (picked.isBefore(_selectedStartDate)) {
          final year = (_months / 12).truncate();
          final month = _months % 12;
          _selectedStartDate = DateTime(
            picked.year - year,
            picked.month - month,
            picked.day,
          );
        }
        _selectedEndDate = picked;
      });
    }
  }

  // Handle duration selection
  void _handleDurationChange(String? value) {
    if (value == null) return;

    setState(() {
      _selectedDuration = value;

      // Update months based on selection
      if (value == '1 month') {
        _months = 1;
      } else if (value == '3 months') {
        _months = 3;
      } else if (value == '6 months') {
        _months = 6;
      } else if (value == '12 months') {
        _months = 12;
      }
      // For custom, we'll use the value from the text field
    });
  }

  // Handle custom months change
  void _handleCustomMonthsChange(String value) {
    if (value.isEmpty) return;

    final parsedValue = int.tryParse(value);
    if (parsedValue != null && parsedValue > 0) {
      setState(() {
        _months = parsedValue;
      });
    }
  }

  Future<void> _submitForm() async {
    if (_formKey.currentState!.validate()) {
      final name = _nameController.text.trim();
      final price = double.parse(_priceController.text);

      // For custom duration, get months from text field
      if (_selectedDuration == 'Custom') {
        final customMonths = int.tryParse(_customMonthsController.text);
        if (customMonths != null && customMonths > 0) {
          _months = customMonths;
        }
      }

      // Validate that if end date is provided, it's not shorter than the subscription recurrence
      if (_selectedEndDate != null) {
        // Calculate the minimum end date based on start date and months
        final minEndDate = DateTime(
          _selectedStartDate.year + (_months ~/ 12),
          _selectedStartDate.month + (_months % 12),
          _selectedStartDate.day,
        );

        if (_selectedEndDate!.isBefore(minEndDate)) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('End date must be at least $_months months after start date'),
              backgroundColor: Colors.red,
              duration: const Duration(seconds: 3),
            ),
          );
          return;
        }
      }

      try {
        // Show loading indicator
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Adding subscription...'),
            duration: Duration(seconds: 1),
          ),
        );

        // Add the subscription using the provider
        await Provider.of<SubscriptionProvider>(
          context,
          listen: false,
        ).addPayment(
          name,
          price,
          _months,
          _selectedStartDate,
          endDate: _selectedEndDate,
          currency: _selectedCurrency,
          labels: _selectedLabels,
        );

        // Show success message
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Subscription added successfully'),
            duration: Duration(seconds: 2),
          ),
        );

        // Close the form
        Navigator.of(context).pop();
      } catch (e) {
        // Show error message
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Error adding subscription: ${e.toString()}'),
            backgroundColor: Colors.red,
            duration: const Duration(seconds: 3),
          ),
        );
      }
    }
  }

  // Build label chips for selection
  List<Widget> _buildLabelChips() {
    final provider = Provider.of<SubscriptionProvider>(context);
    final allLabels = provider.labels;

    return allLabels.map((label) {
      final isSelected = _selectedLabels.any((selectedLabel) => selectedLabel.id == label.id);
      return FilterChip(
        label: Text(label.name),
        selected: isSelected,
        onSelected: (selected) {
          setState(() {
            if (selected) {
              _selectedLabels.add(label);
            } else {
              _selectedLabels.removeWhere((selectedLabel) => selectedLabel.id == label.id);
            }
          });
        },
        backgroundColor: Color(int.parse(label.color.substring(1, 7), radix: 16) + 0xFF000000).withOpacity(0.2),
        selectedColor: Color(int.parse(label.color.substring(1, 7), radix: 16) + 0xFF000000).withOpacity(0.7),
        checkmarkColor: Colors.white,
      );
    }).toList();
  }

  // Show dialog for adding a custom label
  void _showAddLabelDialog() {
    final nameController = TextEditingController();
    final selectedColor = ValueNotifier<Color>(Colors.blue);

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Add Custom Label'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: nameController,
              decoration: const InputDecoration(
                labelText: 'Label Name',
                hintText: 'Enter a name for your label',
              ),
            ),
            const SizedBox(height: 16),
            const Text('Select Color:'),
            const SizedBox(height: 8),
            Wrap(
              spacing: 8,
              children: [
                _buildColorOption(selectedColor, Colors.red),
                _buildColorOption(selectedColor, Colors.pink),
                _buildColorOption(selectedColor, Colors.purple),
                _buildColorOption(selectedColor, Colors.deepPurple),
                _buildColorOption(selectedColor, Colors.indigo),
                _buildColorOption(selectedColor, Colors.blue),
                _buildColorOption(selectedColor, Colors.lightBlue),
                _buildColorOption(selectedColor, Colors.cyan),
                _buildColorOption(selectedColor, Colors.teal),
                _buildColorOption(selectedColor, Colors.green),
                _buildColorOption(selectedColor, Colors.lightGreen),
                _buildColorOption(selectedColor, Colors.lime),
                _buildColorOption(selectedColor, Colors.yellow),
                _buildColorOption(selectedColor, Colors.amber),
                _buildColorOption(selectedColor, Colors.orange),
                _buildColorOption(selectedColor, Colors.deepOrange),
                _buildColorOption(selectedColor, Colors.brown),
                _buildColorOption(selectedColor, Colors.grey),
                _buildColorOption(selectedColor, Colors.blueGrey),
              ],
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () {
              final name = nameController.text.trim();
              if (name.isNotEmpty) {
                final colorHex = '#${selectedColor.value.value.toRadixString(16).substring(2)}';
                Provider.of<SubscriptionProvider>(context, listen: false)
                    .addLabel(name, colorHex)
                    .then((_) {
                  Navigator.of(context).pop();
                });
              }
            },
            child: const Text('Add'),
          ),
        ],
      ),
    );
  }

  // Build a color selection option
  Widget _buildColorOption(ValueNotifier<Color> selectedColor, Color color) {
    return ValueListenableBuilder<Color>(
      valueListenable: selectedColor,
      builder: (context, value, child) {
        final isSelected = value == color;
        return GestureDetector(
          onTap: () => selectedColor.value = color,
          child: Container(
            width: 30,
            height: 30,
            margin: const EdgeInsets.all(2),
            decoration: BoxDecoration(
              color: color,
              shape: BoxShape.circle,
              border: Border.all(
                color: isSelected ? Colors.white : Colors.transparent,
                width: 2,
              ),
              boxShadow: isSelected
                  ? [BoxShadow(color: Colors.black.withOpacity(0.3), blurRadius: 4)]
                  : null,
            ),
          ),
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Add Subscription'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.of(context).pop(),
        ),
      ),
      body: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Form(
            key: _formKey,
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                TextFormField(
                  controller: _nameController,
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
                const SizedBox(height: 16),
                Row(
                  children: [
                    Expanded(
                      flex: 7,
                      child: TextFormField(
                        controller: _priceController,
                        decoration: InputDecoration(
                          labelText: 'Price',
                          hintText: 'Enter the amount',
                          prefixIcon: Icon(_selectedCurrency == Currency.USD.code ? Icons.attach_money : Icons.currency_exchange),
                        ),
                        keyboardType: const TextInputType.numberWithOptions(
                          decimal: true,
                        ),
                        inputFormatters: [
                          FilteringTextInputFormatter.allow(RegExp(r'^\d+\.?\d{0,2}')),
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
                DropdownButtonFormField<String>(
                  value: _selectedDuration,
                  decoration: const InputDecoration(
                    labelText: 'Subscription Recurrence',
                    prefixIcon: Icon(Icons.repeat),
                  ),
                  items: const [
                    DropdownMenuItem(value: '1 month', child: Text('Monthly')),
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
                    DropdownMenuItem(value: 'Custom', child: Text('Custom')),
                  ],
                  onChanged: _handleDurationChange,
                ),
                if (_selectedDuration == 'Custom')
                  Padding(
                    padding: const EdgeInsets.only(top: 16.0),
                    child: TextFormField(
                      controller: _customMonthsController,
                      decoration: const InputDecoration(
                        labelText: 'Custom Months',
                        hintText: 'Enter number of months',
                        prefixIcon: Icon(Icons.calendar_today),
                      ),
                      keyboardType: TextInputType.number,
                      inputFormatters: [FilteringTextInputFormatter.digitsOnly],
                      validator: (value) {
                        if (_selectedDuration == 'Custom') {
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
                      onChanged: _handleCustomMonthsChange,
                    ),
                  ),

                // Date picker
                ListTile(
                  leading: const Icon(Icons.calendar_today),
                  title: const Text('Start Date'),
                  subtitle: Text(
                    '${_selectedStartDate.month}/${_selectedStartDate.day}/${_selectedStartDate.year}',
                  ),
                  onTap: () => _selectStartDate(context),
                ),
                const SizedBox(height: 8),
                ListTile(
                  leading: const Icon(Icons.event_busy),
                  title: const Text('End Date (Optional)'),
                  subtitle: Text(
                    _selectedEndDate != null
                        ? '${_selectedEndDate!.month}/${_selectedEndDate!.day}/${_selectedEndDate!.year}'
                        : 'No end date',
                  ),
                  trailing: _selectedEndDate != null
                      ? IconButton(
                          icon: const Icon(Icons.clear),
                          onPressed: () {
                            setState(() {
                              _selectedEndDate = null;
                            });
                          },
                        )
                      : null,
                  onTap: () => _selectEndDate(context),
                ),
                const SizedBox(height: 16),
                const Divider(),
                const ListTile(
                  title: Text('Labels'),
                  subtitle: Text('Categorize your subscription'),
                  leading: Icon(Icons.label),
                ),
                Wrap(
                  spacing: 8.0,
                  children: _buildLabelChips(),
                ),
                Align(
                  alignment: Alignment.centerRight,
                  child: TextButton.icon(
                    icon: const Icon(Icons.add),
                    label: const Text('Add Custom Label'),
                    onPressed: _showAddLabelDialog,
                  ),
                ),
                const SizedBox(height: 24),
                ElevatedButton(
                  onPressed: _submitForm,
                  style: ElevatedButton.styleFrom(
                    padding: const EdgeInsets.symmetric(vertical: 12),
                  ),
                  child: const Text(
                    'Add Subscription',
                    style: TextStyle(fontSize: 16),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
