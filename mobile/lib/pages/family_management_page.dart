import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../models/family.dart';
import '../models/family_member.dart';
import '../providers/family_provider.dart';
import '../screens/family_member_form_screen.dart';
import '../screens/family_edit_screen.dart';

class FamilyManagementPage extends StatefulWidget {
  const FamilyManagementPage({super.key});

  @override
  State<FamilyManagementPage> createState() => _FamilyManagementPageState();
}

class _FamilyManagementPageState extends State<FamilyManagementPage> {
  late PageController _pageController;

  @override
  void initState() {
    super.initState();
    // Initialize the page controller - we'll set the initial page after we get the data
    _pageController = PageController(viewportFraction: 0.85);
  }

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final familyProvider = Provider.of<FamilyProvider>(context);
    final families = familyProvider.families;
    final selectedFamily = familyProvider.selectedFamily;
    final canEdit = familyProvider.canEditSelectedFamily;

    // Update page controller to show the currently selected family
    if (families.isNotEmpty && selectedFamily != null) {
      final selectedIndex = families.indexWhere(
        (family) => family.id == familyProvider.selectedFamilyId,
      );
      if (selectedIndex != -1 && _pageController.hasClients) {
        // Only animate if the current page is different
        if (_pageController.page?.round() != selectedIndex) {
          _pageController.animateToPage(
            selectedIndex,
            duration: const Duration(milliseconds: 300),
            curve: Curves.easeInOut,
          );
        }
      }
    }

    return SingleChildScrollView(
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: Text(
              'Manage your families',
              style: TextStyle(
                fontSize: 16,
                color: Theme.of(context).textTheme.bodyLarge?.color,
              ),
            ),
          ),
          // Family swipe selector
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
                  Container(
                    height: 100,
                    decoration: BoxDecoration(
                      border: Border.all(color: Colors.grey.shade300),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: PageView.builder(
                      itemCount: families.length,
                      controller: _pageController,
                      onPageChanged: (index) {
                        familyProvider.setSelectedFamilyId(families[index].id);
                      },
                      itemBuilder: (context, index) {
                        final family = families[index];
                        return Padding(
                          padding: const EdgeInsets.all(8.0),
                          child: Card(
                            elevation: 2,
                            child: Center(
                              child: Column(
                                mainAxisAlignment: MainAxisAlignment.center,
                                children: [
                                  Row(
                                    mainAxisAlignment: MainAxisAlignment.center,
                                    children: [
                                      Text(
                                        family.name,
                                        style: const TextStyle(
                                          fontSize: 18,
                                          fontWeight: FontWeight.bold,
                                        ),
                                      ),
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
                                  const SizedBox(height: 4),
                                  Text(
                                    '${family.members.length} members',
                                    style: TextStyle(
                                      color: Colors.grey.shade600,
                                      fontSize: 14,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          ),
                        );
                      },
                    ),
                  ),
                  const SizedBox(height: 8),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      const Icon(Icons.arrow_back_ios, size: 16, color: Colors.grey),
                      const SizedBox(width: 8),
                      Row(
                        mainAxisSize: MainAxisSize.min,
                        children: List.generate(
                          families.length,
                          (index) => Container(
                            width: 8,
                            height: 8,
                            margin: const EdgeInsets.symmetric(horizontal: 2),
                            decoration: BoxDecoration(
                              shape: BoxShape.circle,
                              color: index == families.indexWhere(
                                (family) => family.id == familyProvider.selectedFamilyId,
                              )
                                  ? Theme.of(context).primaryColor
                                  : Colors.grey.shade300,
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(width: 8),
                      const Icon(Icons.arrow_forward_ios, size: 16, color: Colors.grey),
                    ],
                  ),
                  const SizedBox(height: 4),
                  Center(
                    child: Text(
                      'Swipe left or right to change family',
                      style: TextStyle(
                        fontSize: 12,
                        color: Colors.grey.shade600,
                      ),
                    ),
                  ),
                ],
              ),
            ),
          const SizedBox(height: 16),
          // Display family members if a family is selected
          if (selectedFamily != null)
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text(
                        'Family Members:',
                        style: TextStyle(
                          fontWeight: FontWeight.bold,
                          fontSize: 16,
                        ),
                      ),
                      if (selectedFamily != null && canEdit)
                        IconButton(
                          icon: const Icon(Icons.edit),
                          onPressed: () => _navigateToEditFamily(context, selectedFamily.id, familyProvider),
                          tooltip: 'Edit family',
                        ),
                    ],
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
          if (families.isEmpty)
            SizedBox(
              height: MediaQuery.of(context).size.height * 0.6,
              child: _buildEmptyFamiliesView(context),
            ),
          // Add floating action button for creating a new family
          if (!familyProvider.hasOwnedFamily)
            Padding(
              padding: const EdgeInsets.all(16.0),
              child: ElevatedButton.icon(
                onPressed: () => _navigateToAddFamily(context, familyProvider),
                icon: const Icon(Icons.add),
                label: const Text('Create Family'),
                style: ElevatedButton.styleFrom(
                  minimumSize: const Size.fromHeight(50),
                ),
              ),
            ),
        ],
      ),
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

  void _navigateToAddFamily(BuildContext context, FamilyProvider provider) async {
    final result = await Navigator.of(context).push(
      MaterialPageRoute(
        builder: (context) => const FamilyEditScreen(),
      ),
    );
    
    if (result == true) {
      // Family was created, no need to do anything as the provider already updated
    }
  }

  void _navigateToEditFamily(
    BuildContext context,
    String familyId,
    FamilyProvider provider,
  ) async {
    final result = await Navigator.of(context).push(
      MaterialPageRoute(
        builder: (context) => FamilyEditScreen(
          familyId: familyId,
        ),
      ),
    );
    
    if (result == true) {
      // Family was updated, no need to do anything as the provider already updated
    }
  }

  void _showAddMemberDialog(
    BuildContext context,
    String familyId,
    FamilyProvider provider,
  ) async {
    // Navigate to the family member form screen
    final result = await Navigator.of(context).push(
      MaterialPageRoute<Map<String, dynamic>>(
        builder: (context) {
          return FamilyMemberFormScreen(
            familyId: familyId,
          );
        },
      ),
    );

    // Process the result if the user saved the form
    if (result != null) {
      provider.addFamilyMember(
        familyId, 
        result['name'], 
        isKid: result['isKid'],
        email: result['email'],
      );
    }
  }

  void _showEditMemberDialog(
    BuildContext context,
    String familyId,
    FamilyMember member,
    FamilyProvider provider,
  ) async {
    // Navigate to the family member form screen with the existing member data
    final result = await Navigator.of(context).push(
      MaterialPageRoute<Map<String, dynamic>>(
        builder: (context) {
          return FamilyMemberFormScreen(
            familyId: familyId,
            member: member,
          );
        },
      ),
    );

    // Process the result if the user saved the form
    if (result != null) {
      provider.updateFamilyMember(
        familyId,
        member.id,
        result['name'],
        isKid: result['isKid'],
        email: result['email'],
      );
    }
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