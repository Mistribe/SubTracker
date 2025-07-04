import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../models/family_member.dart';
import '../providers/family_member_provider.dart';

class FamilyManagementScreen extends StatelessWidget {
  const FamilyManagementScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final familyMemberProvider = Provider.of<FamilyMemberProvider>(context);
    final familyMembers = familyMemberProvider.familyMembers;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Family Management'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.of(context).pop(),
        ),
      ),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: Text(
              'Manage your family members',
              style: TextStyle(
                fontSize: 16,
                color: Theme.of(context).textTheme.bodyLarge?.color,
              ),
            ),
          ),
          Expanded(
            child: familyMembers.isEmpty
                ? Center(
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(
                          Icons.people_outline,
                          size: 64,
                          color: Theme.of(context).disabledColor,
                        ),
                        const SizedBox(height: 16),
                        Text(
                          'No family members yet',
                          style: TextStyle(
                            fontSize: 18,
                            color: Theme.of(context).disabledColor,
                          ),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          'Add family members to link them to subscriptions',
                          style: TextStyle(
                            fontSize: 14,
                            color: Theme.of(context).disabledColor,
                          ),
                        ),
                      ],
                    ),
                  )
                : ListView.builder(
                    itemCount: familyMembers.length,
                    itemBuilder: (context, index) {
                      final familyMember = familyMembers[index];
                      return ListTile(
                        leading: CircleAvatar(
                          child: Text(
                            familyMember.name.isNotEmpty
                                ? familyMember.name[0].toUpperCase()
                                : '?',
                          ),
                        ),
                        title: Row(
                          children: [
                            Text(familyMember.name),
                            Padding(
                              padding: const EdgeInsets.only(left: 8.0),
                              child: Container(
                                padding: const EdgeInsets.symmetric(
                                  horizontal: 8.0,
                                  vertical: 2.0,
                                ),
                                decoration: BoxDecoration(
                                  color: familyMember.isKid
                                      ? Colors.blue.withValues(alpha: 0.2)
                                      : Colors.green.withValues(alpha: 0.2),
                                  borderRadius: BorderRadius.circular(12.0),
                                ),
                                child: Text(
                                  familyMember.isKid ? 'Kid' : 'Adult',
                                  style: const TextStyle(
                                    fontSize: 12.0,
                                    fontWeight: FontWeight.bold,
                                  ),
                                ),
                              ),
                            ),
                          ],
                        ),
                        trailing: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            IconButton(
                              icon: const Icon(Icons.edit),
                              onPressed: () => _showEditDialog(
                                context,
                                familyMember,
                                familyMemberProvider,
                              ),
                            ),
                            IconButton(
                              icon: const Icon(Icons.delete),
                              onPressed: () => _showDeleteDialog(
                                context,
                                familyMember,
                                familyMemberProvider,
                              ),
                            ),
                          ],
                        ),
                      );
                    },
                  ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showAddDialog(context, familyMemberProvider),
        child: const Icon(Icons.add),
      ),
    );
  }

  void _showAddDialog(BuildContext context, FamilyMemberProvider provider) {
    final nameController = TextEditingController();
    bool isKid = false;

    showDialog(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setState) => AlertDialog(
          title: const Text('Add Family Member'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(
                controller: nameController,
                decoration: const InputDecoration(
                  labelText: 'Name',
                  hintText: 'Enter family member name',
                ),
                textCapitalization: TextCapitalization.words,
                autofocus: true,
              ),
              const SizedBox(height: 16),
              Row(
                children: [
                  Checkbox(
                    value: isKid,
                    onChanged: (value) {
                      setState(() {
                        isKid = value ?? false;
                      });
                    },
                  ),
                  const Text('This is a kid'),
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
                  provider.addFamilyMember(name, isKid: isKid);
                  Navigator.of(context).pop();
                }
              },
              child: const Text('Add'),
            ),
          ],
        ),
      ),
    );
  }

  void _showEditDialog(
    BuildContext context,
    FamilyMember familyMember,
    FamilyMemberProvider provider,
  ) {
    final nameController = TextEditingController(text: familyMember.name);
    bool isKid = familyMember.isKid;

    showDialog(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setState) => AlertDialog(
          title: const Text('Edit Family Member'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(
                controller: nameController,
                decoration: const InputDecoration(
                  labelText: 'Name',
                  hintText: 'Enter family member name',
                ),
                textCapitalization: TextCapitalization.words,
                autofocus: true,
              ),
              const SizedBox(height: 16),
              Row(
                children: [
                  Checkbox(
                    value: isKid,
                    onChanged: (value) {
                      setState(() {
                        isKid = value ?? false;
                      });
                    },
                  ),
                  const Text('This is a kid'),
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
                  provider.updateFamilyMember(
                    familyMember.id,
                    name,
                    isKid: isKid,
                  );
                  Navigator.of(context).pop();
                }
              },
              child: const Text('Save'),
            ),
          ],
        ),
      ),
    );
  }

  void _showDeleteDialog(
    BuildContext context,
    FamilyMember familyMember,
    FamilyMemberProvider provider,
  ) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Delete Family Member'),
        content: Text(
          'Are you sure you want to delete ${familyMember.name}? This action cannot be undone.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () {
              provider.removeFamilyMember(familyMember.id);
              Navigator.of(context).pop();
            },
            style: TextButton.styleFrom(foregroundColor: Colors.red),
            child: const Text('Delete'),
          ),
        ],
      ),
    );
  }
}
