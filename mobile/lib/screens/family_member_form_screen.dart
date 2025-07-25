import 'package:flutter/material.dart';
import '../models/family_member.dart';

class FamilyMemberFormScreen extends StatefulWidget {
  final String familyId;
  final FamilyMember? member; // null for adding new member, non-null for editing
  
  const FamilyMemberFormScreen({
    super.key, 
    required this.familyId,
    this.member,
  });

  @override
  State<FamilyMemberFormScreen> createState() => _FamilyMemberFormScreenState();
}

class _FamilyMemberFormScreenState extends State<FamilyMemberFormScreen> {
  final _nameController = TextEditingController();
  final _emailController = TextEditingController();
  bool _isKid = false;
  
  @override
  void initState() {
    super.initState();
    // If editing an existing member, populate the form
    if (widget.member != null) {
      _nameController.text = widget.member!.name;
      _emailController.text = widget.member?.email ?? '';
      _isKid = widget.member!.isKid;
    }
  }
  
  @override
  void dispose() {
    _nameController.dispose();
    _emailController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.member == null ? 'Add Family Member' : 'Edit Family Member'),
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
                labelText: 'Name',
                hintText: 'Enter family member name',
              ),
              textCapitalization: TextCapitalization.words,
              autofocus: true,
            ),
            const SizedBox(height: 16),
            TextField(
              controller: _emailController,
              decoration: const InputDecoration(
                labelText: 'Email',
                hintText: 'Enter family member email',
              ),
              keyboardType: TextInputType.emailAddress,
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Checkbox(
                  value: _isKid,
                  onChanged: (value) {
                    setState(() {
                      _isKid = value ?? false;
                    });
                  },
                ),
                const Text('This is a kid'),
              ],
            ),
            const SizedBox(height: 32),
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
                    final email = _emailController.text.trim();
                    
                    if (name.isNotEmpty) {
                      final result = {
                        'name': name,
                        'email': email.isNotEmpty ? email : null,
                        'isKid': _isKid,
                      };
                      
                      Navigator.of(context).pop(result);
                    }
                  },
                  child: Text(widget.member == null ? 'Add' : 'Save'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}