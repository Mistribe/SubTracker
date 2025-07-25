import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../models/label.dart';
import '../repositories/label_repository.dart';
import 'package:flutter_colorpicker/flutter_colorpicker.dart';

class LabelManagementPage extends StatefulWidget {
  const LabelManagementPage({super.key});

  @override
  State<LabelManagementPage> createState() => _LabelManagementPageState();
}

class _LabelManagementPageState extends State<LabelManagementPage> {
  late LabelRepository _labelRepository;
  List<Label> _labels = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _labelRepository = Provider.of<LabelRepository>(context, listen: false);
    _loadLabels();
  }

  Future<void> _loadLabels() async {
    setState(() {
      _isLoading = true;
    });

    // Get all labels and sort them (default labels first, then alphabetically)
    final labels = _labelRepository.getAll();
    labels.sort((a, b) {
      if (a.isDefault && !b.isDefault) return -1;
      if (!a.isDefault && b.isDefault) return 1;
      return a.name.compareTo(b.name);
    });

    setState(() {
      _labels = labels;
      _isLoading = false;
    });
  }

  Future<void> _showAddEditLabelDialog({Label? label}) async {
    final isEditing = label != null;
    final nameController = TextEditingController(text: isEditing ? label.name : '');
    Color pickedColor = isEditing ? HexColor.fromHex(label.color) : Colors.blue;
    bool isDefault = isEditing ? label.isDefault : false;

    await showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(isEditing ? 'Edit Label' : 'Add Label'),
        content: SingleChildScrollView(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(
                controller: nameController,
                decoration: const InputDecoration(
                  labelText: 'Label Name',
                  hintText: 'Enter label name',
                ),
                enabled: !isDefault, // Disable editing for default labels
              ),
              const SizedBox(height: 16),
              Text('Select Color:'),
              const SizedBox(height: 8),
              BlockPicker(
                pickerColor: pickedColor,
                onColorChanged: (color) {
                  pickedColor = color;
                },
                availableColors: const [
                  Colors.red,
                  Colors.pink,
                  Colors.purple,
                  Colors.deepPurple,
                  Colors.indigo,
                  Colors.blue,
                  Colors.lightBlue,
                  Colors.cyan,
                  Colors.teal,
                  Colors.green,
                  Colors.lightGreen,
                  Colors.lime,
                  Colors.yellow,
                  Colors.amber,
                  Colors.orange,
                  Colors.deepOrange,
                  Colors.brown,
                  Colors.grey,
                  Colors.blueGrey,
                  Colors.black,
                ],
                itemBuilder: (color, isCurrentColor, onTap) {
                  return Container(
                    margin: const EdgeInsets.all(4),
                    decoration: BoxDecoration(
                      shape: BoxShape.circle,
                      color: color,
                      boxShadow: [
                        if (isCurrentColor)
                          BoxShadow(
                            color: Theme.of(context).colorScheme.shadow.withOpacity(0.4),
                            blurRadius: 4,
                            spreadRadius: 1,
                          )
                      ],
                    ),
                    child: Material(
                      color: Colors.transparent,
                      child: InkWell(
                        onTap: isDefault ? null : onTap, // Disable color change for default labels
                        borderRadius: BorderRadius.circular(50),
                        child: AnimatedOpacity(
                          duration: const Duration(milliseconds: 210),
                          opacity: isCurrentColor ? 1 : 0,
                          child: Icon(
                            Icons.check,
                            color: useWhiteForeground(color) ? Colors.white : Colors.black,
                            size: 24,
                          ),
                        ),
                      ),
                    ),
                  );
                },
              ),
              if (isDefault)
                Padding(
                  padding: const EdgeInsets.only(top: 16.0),
                  child: Text(
                    'Default labels cannot be edited or removed',
                    style: TextStyle(
                      color: Theme.of(context).colorScheme.error,
                      fontStyle: FontStyle.italic,
                    ),
                  ),
                ),
            ],
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Cancel'),
          ),
          if (!isDefault) // Only show Save button for non-default labels
            TextButton(
              onPressed: () async {
                final name = nameController.text.trim();
                if (name.isEmpty) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('Label name cannot be empty')),
                  );
                  return;
                }

                try {
                  if (isEditing) {
                    // Update existing label
                    final updatedLabel = label.copyWith(
                      name: name,
                      color: pickedColor.toHex(),
                      updatedAt: DateTime.now(),
                    );
                    await _labelRepository.update(updatedLabel);
                  } else {
                    // Add new label
                    await _labelRepository.add(name, pickedColor.toHex());
                  }
                } catch (e) {
                  if (mounted) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(content: Text(e.toString())),
                    );
                    return;
                  }
                }

                if (mounted) {
                  Navigator.of(context).pop();
                  _loadLabels(); // Refresh the list
                }
              },
              child: const Text('Save'),
            ),
        ],
      ),
    );
  }

  Future<void> _confirmDeleteLabel(Label label) async {
    if (label.isDefault) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Default labels cannot be removed')),
      );
      return;
    }

    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Delete Label'),
        content: Text('Are you sure you want to delete "${label.name}"?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () => Navigator.of(context).pop(true),
            child: const Text('Delete'),
          ),
        ],
      ),
    );

    if (confirmed == true && mounted) {
      try {
        await _labelRepository.delete(label.id);
        _loadLabels(); // Refresh the list
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text(e.toString())),
          );
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return _isLoading
        ? const Center(child: CircularProgressIndicator())
        : Column(
            children: [
              Padding(
                padding: const EdgeInsets.all(16.0),
                child: Text(
                  'Manage your labels',
                  style: TextStyle(
                    fontSize: 16,
                    color: Theme.of(context).textTheme.bodyLarge?.color,
                  ),
                ),
              ),
              Expanded(
                child: ListView.builder(
                  itemCount: _labels.length,
                  itemBuilder: (context, index) {
                    final label = _labels[index];
                    return ListTile(
                      leading: CircleAvatar(
                        backgroundColor: HexColor.fromHex(label.color),
                        child: label.isDefault
                            ? Icon(
                                Icons.lock, 
                                color: useWhiteForeground(HexColor.fromHex(label.color)) 
                                    ? Colors.white 
                                    : Colors.black, 
                                size: 16
                              )
                            : null,
                      ),
                      title: Text(label.name),
                      subtitle: label.isDefault
                          ? const Text('Default label', style: TextStyle(fontStyle: FontStyle.italic))
                          : const Text('Custom label'),
                      trailing: Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          IconButton(
                            icon: const Icon(Icons.edit),
                            onPressed: label.isDefault
                                ? null // Disable for default labels
                                : () => _showAddEditLabelDialog(label: label),
                            tooltip: label.isDefault ? 'Default labels cannot be edited' : 'Edit label',
                          ),
                          IconButton(
                            icon: const Icon(Icons.delete),
                            onPressed: label.isDefault
                                ? null // Disable for default labels
                                : () => _confirmDeleteLabel(label),
                            tooltip: label.isDefault ? 'Default labels cannot be removed' : 'Delete label',
                          ),
                        ],
                      ),
                    );
                  },
                ),
              ),
              Padding(
                padding: const EdgeInsets.all(16.0),
                child: ElevatedButton.icon(
                  onPressed: () => _showAddEditLabelDialog(),
                  icon: const Icon(Icons.add),
                  label: const Text('Add Label'),
                  style: ElevatedButton.styleFrom(
                    minimumSize: const Size.fromHeight(50),
                  ).copyWith(
                    backgroundColor: WidgetStatePropertyAll(Theme.of(context).colorScheme.primaryContainer),
                    foregroundColor: WidgetStatePropertyAll(Theme.of(context).colorScheme.onPrimaryContainer),
                  ),
                ),
              ),
            ],
          );
  }
}

// Extension to convert between Color and hex string
extension HexColor on Color {
  static Color fromHex(String hexString) {
    final buffer = StringBuffer();
    if (hexString.length == 6 || hexString.length == 7) buffer.write('ff');
    buffer.write(hexString.replaceFirst('#', ''));
    return Color(int.parse(buffer.toString(), radix: 16));
  }

  String toHex({bool leadingHashSign = true}) => '${leadingHashSign ? '#' : ''}'
      '${alpha.toRadixString(16).padLeft(2, '0')}'
      '${red.toRadixString(16).padLeft(2, '0')}'
      '${green.toRadixString(16).padLeft(2, '0')}'
      '${blue.toRadixString(16).padLeft(2, '0')}';
}
