import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/authentication_provider.dart';
import '../providers/sync_provider.dart';

class AccountScreen extends StatelessWidget {
  const AccountScreen({super.key});

  // Format DateTime to a readable string
  String _formatDateTime(DateTime dateTime) {
    final now = DateTime.now();
    final difference = now.difference(dateTime);

    if (difference.inDays > 0) {
      return '${difference.inDays} day${difference.inDays == 1 ? '' : 's'} ago';
    } else if (difference.inHours > 0) {
      return '${difference.inHours} hour${difference.inHours == 1 ? '' : 's'} ago';
    } else if (difference.inMinutes > 0) {
      return '${difference.inMinutes} minute${difference.inMinutes == 1 ? '' : 's'} ago';
    } else {
      return 'Just now';
    }
  }

  @override
  Widget build(BuildContext context) {
    final authProvider = Provider.of<AuthenticationProvider>(context);
    final syncProvider = Provider.of<SyncProvider>(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Account'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.of(context).pop(),
        ),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // User profile section
            const Text(
              'Profile Information',
              style: TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 16),
            Card(
              elevation: 2,
              child: Padding(
                padding: const EdgeInsets.all(16.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    ListTile(
                      leading: const CircleAvatar(
                        child: Icon(Icons.person),
                      ),
                      title: Text(
                        authProvider.user?.displayName ?? 'User',
                        style: const TextStyle(
                          fontWeight: FontWeight.bold,
                          fontSize: 18,
                        ),
                      ),
                      subtitle: Text(
                        authProvider.user?.email ?? 'email',
                      ),
                    ),
                    const SizedBox(height: 8),
                    const Text(
                      'Your account is managed through our authentication provider. To update your profile information, please visit the account settings on our authentication provider\'s website.',
                      style: TextStyle(
                        color: Colors.grey,
                        fontSize: 12,
                      ),
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 32),

            // Synchronization section
            const Text(
              'Synchronization',
              style: TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 16),
            FutureBuilder<int>(
              future: syncProvider.getPendingOperationsCount(),
              builder: (context, snapshot) {
                final pendingCount = snapshot.data ?? 0;
                final isFullySynced = !syncProvider.hasPendingOperations;

                return Card(
                  elevation: 2,
                  child: Padding(
                    padding: const EdgeInsets.all(16.0),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        // Last sync time
                        Row(
                          children: [
                            const Icon(Icons.access_time, color: Colors.grey),
                            const SizedBox(width: 8),
                            Expanded(
                              child: Text(
                                'Last sync: ${syncProvider.lastSyncTime != null ? _formatDateTime(syncProvider.lastSyncTime!) : 'Never'}',
                                style: const TextStyle(fontSize: 14),
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 12),

                        // Sync status
                        Row(
                          children: [
                            Icon(
                              isFullySynced ? Icons.check_circle : Icons.sync_problem,
                              color: isFullySynced ? Colors.green : Colors.orange,
                            ),
                            const SizedBox(width: 8),
                            Expanded(
                              child: Text(
                                isFullySynced 
                                  ? 'Fully synchronized' 
                                  : 'Synchronization pending',
                                style: TextStyle(
                                  fontSize: 14,
                                  color: isFullySynced ? Colors.green : Colors.orange,
                                  fontWeight: FontWeight.bold,
                                ),
                              ),
                            ),
                          ],
                        ),

                        // Pending operations count
                        if (pendingCount > 0) ...[
                          const SizedBox(height: 12),
                          Row(
                            children: [
                              const Icon(Icons.pending_actions, color: Colors.orange),
                              const SizedBox(width: 8),
                              Text(
                                '$pendingCount package${pendingCount == 1 ? '' : 's'} pending synchronization',
                                style: const TextStyle(fontSize: 14),
                              ),
                            ],
                          ),
                        ],

                        const SizedBox(height: 16),

                        // Manual sync button
                        ElevatedButton.icon(
                          onPressed: syncProvider.isSyncing 
                            ? null 
                            : () => syncProvider.sync(),
                          icon: syncProvider.isSyncing 
                            ? const SizedBox(
                                width: 20,
                                height: 20,
                                child: CircularProgressIndicator(strokeWidth: 2),
                              )
                            : const Icon(Icons.sync),
                          label: Text(
                            syncProvider.isSyncing 
                              ? 'Syncing...' 
                              : 'Sync Now'
                          ),
                          style: ElevatedButton.styleFrom(
                            minimumSize: const Size(double.infinity, 40),
                          ),
                        ),
                      ],
                    ),
                  ),
                );
              },
            ),

            const SizedBox(height: 32),

            // Account actions section
            const Text(
              'Account Actions',
              style: TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 16),

            // Sign out button
            ElevatedButton.icon(
              onPressed: () async {
                await authProvider.signOut();
                syncProvider.updateSyncEnabled(isAuthenticated: false);
                if (context.mounted) {
                  Navigator.of(context).pop();
                }
              },
              icon: const Icon(Icons.logout),
              label: const Text('Sign Out'),
              style: ElevatedButton.styleFrom(
                minimumSize: const Size(double.infinity, 50),
              ),
            ),
            const SizedBox(height: 16),

            // Delete account button (currently disabled with explanation)
            ElevatedButton.icon(
              onPressed: () {
                showDialog(
                  context: context,
                  builder: (context) => AlertDialog(
                    title: const Text('Account Deletion'),
                    content: const Text(
                      'To delete your account, please visit our authentication provider\'s website and follow their account deletion process. This will ensure all your data is properly removed from our systems.',
                    ),
                    actions: [
                      TextButton(
                        onPressed: () => Navigator.of(context).pop(),
                        child: const Text('OK'),
                      ),
                    ],
                  ),
                );
              },
              icon: const Icon(Icons.delete_forever),
              label: const Text('Delete Account'),
              style: ElevatedButton.styleFrom(
                minimumSize: const Size(double.infinity, 50),
                backgroundColor: Colors.red,
                foregroundColor: Colors.white,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
