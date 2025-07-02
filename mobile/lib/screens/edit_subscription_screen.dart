import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import '../providers/subscription_provider.dart';
import '../models/subscription.dart';
import '../models/label.dart';
import '../models/currency.dart';

class EditSubscriptionScreen extends StatefulWidget {
  final Subscription subscription;

  const EditSubscriptionScreen({super.key, required this.subscription});

  @override
  State<EditSubscriptionScreen> createState() => _EditSubscriptionScreenState();
}

class _EditSubscriptionScreenState extends State<EditSubscriptionScreen> {
  final _formKey = GlobalKey<FormState>();
  late final TextEditingController _nameController;
  late final TextEditingController _priceController;
  late final TextEditingController _customMonthsController;
  late int _months;
  String _selectedDuration = '1 month'; // Default value
  late DateTime _startDate;
  DateTime? _endDate;
  late String _selectedCurrency;
  late List<Label> _selectedLabels;

  // Get currencies from the Currency enum
  final List<String> _currencies = Currency.codes;

  @override
  void initState() {
    super.initState();
    // Initialize controllers with existing subscription data
    _nameController = TextEditingController(text: widget.subscription.name);

    // Get the current subscription detail
    final currentDetail = widget.subscription.getLastPaymentDetail();
    _priceController = TextEditingController(
      text: currentDetail.price.toString(),
    );
    _customMonthsController = TextEditingController();

    // Initialize months from subscription detail
    _months = currentDetail.months;
    _startDate = currentDetail.startDate;
    _endDate = currentDetail.endDate;
    _selectedCurrency = currentDetail.currency;
    _selectedLabels = List<Label>.from(widget.subscription.labels);

    // Set the selected duration based on months
    if (_months == 1) {
      _selectedDuration = '1 month';
    } else if (_months == 3) {
      _selectedDuration = '3 months';
    } else if (_months == 6) {
      _selectedDuration = '6 months';
    } else if (_months == 12) {
      _selectedDuration = '12 months';
    } else {
      _selectedDuration = 'Custom';
      _customMonthsController.text = _months.toString();
    }
  }

  @override
  void dispose() {
    _nameController.dispose();
    _priceController.dispose();
    _customMonthsController.dispose();
    super.dispose();
  }

  // Date picker for start date
  Future<void> _selectStartDate(BuildContext context) async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate: _startDate,
      firstDate: DateTime(2000),
      lastDate: DateTime(2101),
    );
    if (picked != null && picked != _startDate) {
      setState(() {
        _startDate = picked;
      });
    }
  }

  // Date picker for end date
  Future<void> _selectEndDate(BuildContext context) async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate: _endDate ?? DateTime.now(),
      firstDate: _startDate,
      lastDate: DateTime(2101),
    );
    if (picked != null && picked != _endDate) {
      setState(() {
        _endDate = picked;
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
      final price = double.tryParse(_priceController.text) ?? 0.0;

      // For custom duration, get months from text field
      if (_selectedDuration == 'Custom') {
        final customMonths = int.tryParse(_customMonthsController.text);
        if (customMonths != null && customMonths > 0) {
          _months = customMonths;
        }
      }

      // Validate that if end date is provided, it's not shorter than the subscription recurrence
      if (_endDate != null) {
        // Calculate the minimum end date based on start date and months
        final minEndDate = DateTime(
          _startDate.year + (_months ~/ 12),
          _startDate.month + (_months % 12),
          _startDate.day,
        );

        if (_endDate!.isBefore(minEndDate)) {
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
            content: Text('Updating subscription...'),
            duration: Duration(seconds: 1),
          ),
        );

        // Update the subscription name and labels
        await Provider.of<SubscriptionProvider>(
          context,
          listen: false,
        ).updateSubscription(widget.subscription.id, name);

        // Update the subscription labels
        await Provider.of<SubscriptionProvider>(
          context,
          listen: false,
        ).updateSubscriptionLabels(widget.subscription.id, _selectedLabels);

        // Update the subscription detail
        final currentDetail = widget.subscription.getLastPaymentDetail();
        await Provider.of<SubscriptionProvider>(
          context,
          listen: false,
        ).updateSubscriptionPayment(
          widget.subscription.id,
          currentDetail.id,
          price,
          _months,
          _startDate,
          endDate: _endDate,
          currency: _selectedCurrency,
        );

        // Show success message
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Payment updated successfully'),
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
        title: const Text(
          'Edit Subscription',
          style: TextStyle(
            fontSize: 20,
            fontWeight: FontWeight.w600,
          ),
        ),
        elevation: 0,
      ),
      body: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Form(
            key: _formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                // Basic Information Section Title
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

                // Basic Information Card
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
                      ],
                    ),
                  ),
                ),

                if (widget.subscription.isActive) ...[
                  const SizedBox(height: 24),

                  // Subscription Details Section Title
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

                  // Subscription Details Card
                  Card(
                    elevation: 2,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(16),
                    ),
                    child: Padding(
                      padding: const EdgeInsets.all(16.0),
                      child: Column(
                        children: [
                          Row(
                            children: [
                              Expanded(
                                flex: 7,
                                child: TextFormField(
                                  controller: _priceController,
                                  decoration: InputDecoration(
                                    labelText: 'Price',
                                    hintText: 'Enter the price',
                                    prefixIcon: Icon(_selectedCurrency == Currency.USD.code ? Icons.attach_money : Icons.currency_exchange),
                                  ),
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
                                    final price = double.tryParse(value);
                                    if (price == null || price <= 0) {
                                      return 'Please enter a valid price';
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
                                inputFormatters: [
                                  FilteringTextInputFormatter.digitsOnly,
                                ],
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
                          const SizedBox(height: 16),

                          // Date pickers with modern styling
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
                                color: Colors.blue.withOpacity(Theme.of(context).brightness == Brightness.dark ? 0.8 : 1.0),
                              ),
                              title: const Text('Start Date'),
                              subtitle: Text(
                                '${_startDate.month}/${_startDate.day}/${_startDate.year}',
                              ),
                              onTap: () => _selectStartDate(context),
                            ),
                          ),

                          const SizedBox(height: 16),

                          Container(
                            decoration: BoxDecoration(
                              borderRadius: BorderRadius.circular(12),
                              border: Border.all(
                                color: Theme.of(context).dividerColor.withOpacity(0.5),
                              ),
                            ),
                            child: ListTile(
                              leading: Icon(
                                Icons.event_busy,
                                color: Colors.orange.withOpacity(Theme.of(context).brightness == Brightness.dark ? 0.8 : 1.0),
                              ),
                              title: const Text('End Date (Optional)'),
                              subtitle: Text(
                                _endDate != null
                                    ? '${_endDate!.month}/${_endDate!.day}/${_endDate!.year}'
                                    : 'No end date',
                              ),
                              trailing: _endDate != null
                                  ? IconButton(
                                      icon: const Icon(Icons.clear),
                                      onPressed: () {
                                        setState(() {
                                          _endDate = null;
                                        });
                                      },
                                    )
                                  : null,
                              onTap: () => _selectEndDate(context),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ],

                const SizedBox(height: 24),

                // Labels Section Title
                Padding(
                  padding: const EdgeInsets.only(left: 4, bottom: 12),
                  child: Text(
                    'Labels',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.w600,
                      color: Theme.of(context).colorScheme.primary,
                    ),
                  ),
                ),

                // Labels Card
                Card(
                  elevation: 2,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(16),
                  ),
                  child: Padding(
                    padding: const EdgeInsets.all(16.0),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
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
                      ],
                    ),
                  ),
                ),

                const SizedBox(height: 24),

                // Submit Button
                ElevatedButton(
                  onPressed: _submitForm,
                  style: ElevatedButton.styleFrom(
                    padding: const EdgeInsets.symmetric(vertical: 14),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                  child: const Text(
                    'Update Subscription',
                    style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
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
