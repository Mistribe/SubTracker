import 'package:flutter/material.dart';
import '../models/family.dart';
import '../providers/family_provider.dart';
import 'package:provider/provider.dart';

class FamilyEditScreen extends StatefulWidget {
  final String? familyId; // null for adding new family, non-null for editing

  const FamilyEditScreen({super.key, this.familyId});

  @override
  State<FamilyEditScreen> createState() => _FamilyEditScreenState();
}

class _FamilyEditScreenState extends State<FamilyEditScreen> {
  final _nameController = TextEditingController();
  bool _haveJointAccount = true;

  @override
  void initState() {
    super.initState();
    // If editing an existing family, populate the form
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (widget.familyId != null) {
        final familyProvider = Provider.of<FamilyProvider>(
          context,
          listen: false,
        );
        final family = familyProvider.getFamilyById(widget.familyId!);
        if (family != null) {
          setState(() {
            _nameController.text = family.name;
            _haveJointAccount = family.haveJointAccount;
          });
        }
      }
    });
  }

  @override
  void dispose() {
    _nameController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final familyProvider = Provider.of<FamilyProvider>(context);
    final family = widget.familyId != null
        ? familyProvider.getFamilyById(widget.familyId!)
        : null;

    return Scaffold(
      appBar: AppBar(
        title: Text(widget.familyId == null ? 'Create Family' : 'Edit Family'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.of(context).pop(),
        ),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            TextField(
              controller: _nameController,
              decoration: const InputDecoration(
                labelText: 'Family Name',
                hintText: 'Enter family name',
              ),
              textCapitalization: TextCapitalization.words,
              autofocus: true,
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Checkbox(
                  value: _haveJointAccount,
                  onChanged: (value) {
                    setState(() {
                      _haveJointAccount = value ?? true;
                    });
                  },
                ),
                const Text('Have joint account'),
              ],
            ),
            const SizedBox(height: 24),
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                TextButton(
                  onPressed: () => Navigator.of(context).pop(),
                  child: const Text('Cancel'),
                ),
                const SizedBox(width: 16),
                ElevatedButton(
                  onPressed: () {
                    final name = _nameController.text.trim();
                    if (name.isNotEmpty) {
                      if (widget.familyId == null) {
                        // Create new family
                        familyProvider.createFamily(
                          name,
                          haveJointAccount: _haveJointAccount,
                        );
                      } else {
                        // Update existing family
                        familyProvider.updateFamily(
                          widget.familyId!,
                          name,
                          haveJointAccount: _haveJointAccount,
                        );
                      }
                      Navigator.of(context).pop(true); // Return success
                    }
                  },
                  child: Text(widget.familyId == null ? 'Create' : 'Save'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

}
