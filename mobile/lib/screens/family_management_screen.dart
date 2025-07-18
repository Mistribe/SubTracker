import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../models/family.dart';
import '../models/family_member.dart';
import '../providers/family_provider.dart';

class FamilyManagementScreen extends StatefulWidget {
  const FamilyManagementScreen({super.key});

  @override
  State<FamilyManagementScreen> createState() => _FamilyManagementScreenState();
}

class _FamilyManagementScreenState extends State<FamilyManagementScreen> {
  bool _isEditMode = false;
  final _nameController = TextEditingController();
  bool _haveJointAccount = true;

  @override
  void dispose() {
    _nameController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final familyProvider = Provider.of<FamilyProvider>(context);
    final families = familyProvider.families;
    final selectedFamily = familyProvider.selectedFamily;
    final canEdit = familyProvider.canEditSelectedFamily;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Family Management'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.of(context).pop(),
        ),
        actions: [
          if (selectedFamily != null && canEdit)
            IconButton(
              icon: Icon(_isEditMode ? Icons.visibility : Icons.edit),
              onPressed: () {
                setState(() {
                  _isEditMode = !_isEditMode;
                  if (_isEditMode) {
                    // Initialize form with current values
                    _nameController.text = selectedFamily.name;
                    _haveJointAccount = selectedFamily.haveJointAccount;
                  }
                });
              },
              tooltip: _isEditMode ? 'View mode' : 'Edit mode',
            ),
        ],
      ),
      body: SingleChildScrollView(
        child: Column(
          children: [
            Padding(
              padding: const EdgeInsets.all(16.0),
              child: Text(
                _isEditMode ? 'Edit your family' : 'Manage your families',
                style: TextStyle(
                  fontSize: 16,
                  color: Theme.of(context).textTheme.bodyLarge?.color,
                ),
              ),
            ),
            // Family dropdown selector
            if (families.isNotEmpty)
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Active Family:',
                      style: TextStyle(
                        fontWeight: FontWeight.bold,
                        fontSize: 16,
                      ),
                    ),
                    const SizedBox(height: 8),
                    DropdownButtonFormField<String>(
                      value: familyProvider.selectedFamilyId,
                      decoration: const InputDecoration(
                        border: OutlineInputBorder(),
                        contentPadding: EdgeInsets.symmetric(
                          horizontal: 16,
                          vertical: 8,
                        ),
                      ),
                      items: families.map((family) {
                        return DropdownMenuItem<String>(
                          value: family.id,
                          child: Row(
                            children: [
                              Text(family.name),
                              if (family.isOwner)
                                Padding(
                                  padding: const EdgeInsets.only(left: 8.0),
                                  child: Container(
                                    padding: const EdgeInsets.symmetric(
                                      horizontal: 8.0,
                                      vertical: 2.0,
                                    ),
                                    decoration: BoxDecoration(
                                      color: Colors.green.withOpacity(0.2),
                                      borderRadius: BorderRadius.circular(12.0),
                                    ),
                                    child: const Text(
                                      'Owner',
                                      style: TextStyle(
                                        fontSize: 12.0,
                                        fontWeight: FontWeight.bold,
                                      ),
                                    ),
                                  ),
                                ),
                            ],
                          ),
                        );
                      }).toList(),
                      onChanged: (String? newValue) {
                        if (newValue != null) {
                          familyProvider.setSelectedFamilyId(newValue);
                          setState(() {
                            _isEditMode = false;
                          });
                        }
                      },
                    ),
                  ],
                ),
              ),
            const SizedBox(height: 16),
            // Display family members if a family is selected
            if (selectedFamily != null && !_isEditMode)
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Family Members:',
                      style: TextStyle(
                        fontWeight: FontWeight.bold,
                        fontSize: 16,
                      ),
                    ),
                    const SizedBox(height: 8),
                    selectedFamily.members.isEmpty
                        ? const Padding(
                            padding: EdgeInsets.symmetric(vertical: 8.0),
                            child: Text('No family members yet'),
                          )
                        : ListView.builder(
                            shrinkWrap: true,
                            physics: const NeverScrollableScrollPhysics(),
                            itemCount: selectedFamily.members.length,
                            itemBuilder: (context, index) {
                              final member = selectedFamily.members[index];
                              return ListTile(
                                leading: CircleAvatar(
                                  child: Text(
                                    member.name.isNotEmpty
                                        ? member.name[0].toUpperCase()
                                        : '?',
                                  ),
                                ),
                                title: Row(
                                  children: [
                                    Text(member.name),
                                    Padding(
                                      padding: const EdgeInsets.only(left: 8.0),
                                      child: Container(
                                        padding: const EdgeInsets.symmetric(
                                          horizontal: 8.0,
                                          vertical: 2.0,
                                        ),
                                        decoration: BoxDecoration(
                                          color: member.isKid
                                              ? Colors.blue.withOpacity(0.2)
                                              : Colors.green.withOpacity(0.2),
                                          borderRadius: BorderRadius.circular(
                                            12.0,
                                          ),
                                        ),
                                        child: Text(
                                          member.isKid ? 'Kid' : 'Adult',
                                          style: const TextStyle(
                                            fontSize: 12.0,
                                            fontWeight: FontWeight.bold,
                                          ),
                                        ),
                                      ),
                                    ),
                                  ],
                                ),
                                trailing: canEdit
                                    ? Row(
                                        mainAxisSize: MainAxisSize.min,
                                        children: [
                                          IconButton(
                                            icon: const Icon(Icons.edit),
                                            onPressed: () =>
                                                _showEditMemberDialog(
                                                  context,
                                                  selectedFamily.id,
                                                  member,
                                                  familyProvider,
                                                ),
                                          ),
                                          IconButton(
                                            icon: const Icon(Icons.delete),
                                            onPressed: () =>
                                                _showDeleteMemberDialog(
                                                  context,
                                                  selectedFamily.id,
                                                  member,
                                                  familyProvider,
                                                ),
                                          ),
                                        ],
                                      )
                                    : null,
                              );
                            },
                          ),
                    if (canEdit)
                      Padding(
                        padding: const EdgeInsets.only(top: 16.0, bottom: 16.0),
                        child: ElevatedButton.icon(
                          onPressed: () => _showAddMemberDialog(
                            context,
                            selectedFamily.id,
                            familyProvider,
                          ),
                          icon: const Icon(Icons.add),
                          label: const Text('Add Family Member'),
                        ),
                      ),
                  ],
                ),
              ),
            if (selectedFamily != null && _isEditMode)
              _buildFamilyEditForm(context, selectedFamily, familyProvider)
            else if (families.isEmpty)
              SizedBox(
                height: MediaQuery.of(context).size.height * 0.6,
                child: _buildEmptyFamiliesView(context),
              ),
          ],
        ),
      ),
      floatingActionButton: !_isEditMode && !familyProvider.hasOwnedFamily
          ? FloatingActionButton(
              onPressed: () => _showAddFamilyDialog(context, familyProvider),
              child: const Icon(Icons.add),
            )
          : null,
    );
  }

  Widget _buildEmptyFamiliesView(BuildContext context) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.family_restroom,
            size: 64,
            color: Theme.of(context).disabledColor,
          ),
          const SizedBox(height: 16),
          Text(
            'No families yet',
            style: TextStyle(
              fontSize: 18,
              color: Theme.of(context).disabledColor,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'Create a family to manage family members',
            style: TextStyle(
              fontSize: 14,
              color: Theme.of(context).disabledColor,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildFamiliesList(
    BuildContext context,
    List<Family> families,
    FamilyProvider provider,
  ) {
    return ListView.builder(
      itemCount: families.length,
      itemBuilder: (context, index) {
        final family = families[index];
        return ListTile(
          leading: CircleAvatar(
            child: Text(
              family.name.isNotEmpty ? family.name[0].toUpperCase() : '?',
            ),
          ),
          title: Row(
            children: [
              Text(family.name),
              if (family.isOwner)
                Padding(
                  padding: const EdgeInsets.only(left: 8.0),
                  child: Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 8.0,
                      vertical: 2.0,
                    ),
                    decoration: BoxDecoration(
                      color: Colors.green.withOpacity(0.2),
                      borderRadius: BorderRadius.circular(12.0),
                    ),
                    child: const Text(
                      'Owner',
                      style: TextStyle(
                        fontSize: 12.0,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ),
            ],
          ),
          subtitle: Text('${family.members.length} members'),
          selected: provider.selectedFamilyId == family.id,
          onTap: () {
            provider.setSelectedFamilyId(family.id);
            setState(() {
              _isEditMode = false;
            });
          },
          trailing: family.isOwner
              ? Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    IconButton(
                      icon: const Icon(Icons.edit),
                      onPressed: () {
                        provider.setSelectedFamilyId(family.id);
                        setState(() {
                          _isEditMode = true;
                          _nameController.text = family.name;
                          _haveJointAccount = family.haveJointAccount;
                        });
                      },
                    ),
                    IconButton(
                      icon: const Icon(Icons.delete),
                      onPressed: () =>
                          _showDeleteFamilyDialog(context, family, provider),
                    ),
                  ],
                )
              : null,
        );
      },
    );
  }

  Widget _buildFamilyEditForm(
    BuildContext context,
    Family family,
    FamilyProvider provider,
  ) {
    return Expanded(
      child: SingleChildScrollView(
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
            const Text(
              'Family Members',
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 8),
            if (family.members.isEmpty)
              const Padding(
                padding: EdgeInsets.symmetric(vertical: 16.0),
                child: Text('No family members yet'),
              )
            else
              ListView.builder(
                shrinkWrap: true,
                physics: const NeverScrollableScrollPhysics(),
                itemCount: family.members.length,
                itemBuilder: (context, index) {
                  final member = family.members[index];
                  return ListTile(
                    leading: CircleAvatar(
                      child: Text(
                        member.name.isNotEmpty
                            ? member.name[0].toUpperCase()
                            : '?',
                      ),
                    ),
                    title: Row(
                      children: [
                        Text(member.name),
                        Padding(
                          padding: const EdgeInsets.only(left: 8.0),
                          child: Container(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 8.0,
                              vertical: 2.0,
                            ),
                            decoration: BoxDecoration(
                              color: member.isKid
                                  ? Colors.blue.withOpacity(0.2)
                                  : Colors.green.withOpacity(0.2),
                              borderRadius: BorderRadius.circular(12.0),
                            ),
                            child: Text(
                              member.isKid ? 'Kid' : 'Adult',
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
                          onPressed: () => _showEditMemberDialog(
                            context,
                            family.id,
                            member,
                            provider,
                          ),
                        ),
                        IconButton(
                          icon: const Icon(Icons.delete),
                          onPressed: () => _showDeleteMemberDialog(
                            context,
                            family.id,
                            member,
                            provider,
                          ),
                        ),
                      ],
                    ),
                  );
                },
              ),
            const SizedBox(height: 16),
            ElevatedButton.icon(
              onPressed: () =>
                  _showAddMemberDialog(context, family.id, provider),
              icon: const Icon(Icons.add),
              label: const Text('Add Family Member'),
            ),
            const SizedBox(height: 24),
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                TextButton(
                  onPressed: () {
                    setState(() {
                      _isEditMode = false;
                    });
                  },
                  child: const Text('Cancel'),
                ),
                const SizedBox(width: 16),
                ElevatedButton(
                  onPressed: () {
                    final name = _nameController.text.trim();
                    if (name.isNotEmpty) {
                      provider.updateFamily(
                        family.id,
                        name,
                        haveJointAccount: _haveJointAccount,
                      );
                      setState(() {
                        _isEditMode = false;
                      });
                    }
                  },
                  child: const Text('Save'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  void _showAddFamilyDialog(BuildContext context, FamilyProvider provider) {
    final nameController = TextEditingController();
    bool haveJointAccount = true;

    showDialog(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setState) => AlertDialog(
          title: const Text('Create Family'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(
                controller: nameController,
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
                    value: haveJointAccount,
                    onChanged: (value) {
                      setState(() {
                        haveJointAccount = value ?? true;
                      });
                    },
                  ),
                  const Text('Have joint account'),
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
                  provider.createFamily(
                    name,
                    haveJointAccount: haveJointAccount,
                  );
                  Navigator.of(context).pop();
                }
              },
              child: const Text('Create'),
            ),
          ],
        ),
      ),
    );
  }

  void _showDeleteFamilyDialog(
    BuildContext context,
    Family family,
    FamilyProvider provider,
  ) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Delete Family'),
        content: Text(
          'Are you sure you want to delete ${family.name}? This action cannot be undone.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () {
              provider.deleteFamily(family.id);
              Navigator.of(context).pop();
            },
            style: TextButton.styleFrom(foregroundColor: Colors.red),
            child: const Text('Delete'),
          ),
        ],
      ),
    );
  }

  void _showAddMemberDialog(
    BuildContext context,
    String familyId,
    FamilyProvider provider,
  ) {
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
                  provider.addFamilyMember(familyId, name, isKid: isKid);
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

  void _showEditMemberDialog(
    BuildContext context,
    String familyId,
    FamilyMember member,
    FamilyProvider provider,
  ) {
    final nameController = TextEditingController(text: member.name);
    bool isKid = member.isKid;

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
                    familyId,
                    member.id,
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

  void _showDeleteMemberDialog(
    BuildContext context,
    String familyId,
    FamilyMember member,
    FamilyProvider provider,
  ) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Delete Family Member'),
        content: Text(
          'Are you sure you want to delete ${member.name}? This action cannot be undone.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () {
              provider.removeFamilyMember(familyId, member.id);
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
