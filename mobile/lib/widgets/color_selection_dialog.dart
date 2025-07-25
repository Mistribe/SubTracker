import 'package:flutter/material.dart';

class ColorSelectionDialog extends StatelessWidget {
  final TextEditingController nameController;
  final ValueNotifier<Color> selectedColor;
  final Function(String, String) onAddLabel;

  const ColorSelectionDialog({
    super.key,
    required this.nameController,
    required this.selectedColor,
    required this.onAddLabel,
  });

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
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
              _buildColorOption(Colors.red),
              _buildColorOption(Colors.pink),
              _buildColorOption(Colors.purple),
              _buildColorOption(Colors.deepPurple),
              _buildColorOption(Colors.indigo),
              _buildColorOption(Colors.blue),
              _buildColorOption(Colors.lightBlue),
              _buildColorOption(Colors.cyan),
              _buildColorOption(Colors.teal),
              _buildColorOption(Colors.green),
              _buildColorOption(Colors.lightGreen),
              _buildColorOption(Colors.lime),
              _buildColorOption(Colors.yellow),
              _buildColorOption(Colors.amber),
              _buildColorOption(Colors.orange),
              _buildColorOption(Colors.deepOrange),
              _buildColorOption(Colors.brown),
              _buildColorOption(Colors.grey),
              _buildColorOption(Colors.blueGrey),
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
              final colorHex = '#${selectedColor.value.value.toRadixString(16).padLeft(8, '0').substring(2)}';
              onAddLabel(name, colorHex);
              Navigator.of(context).pop();
            }
          },
          child: const Text('Add'),
        ),
      ],
    );
  }

  Widget _buildColorOption(Color color) {
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
                  ? [
                      BoxShadow(
                        color: Colors.black.withOpacity(0.3),
                        blurRadius: 4,
                      ),
                    ]
                  : null,
            ),
          ),
        );
      },
    );
  }
}
