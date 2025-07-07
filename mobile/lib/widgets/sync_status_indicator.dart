import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/sync_provider.dart';

/// A widget that displays the current synchronization status
class SyncStatusIndicator extends StatelessWidget {
  const SyncStatusIndicator({super.key});

  @override
  Widget build(BuildContext context) {
    return Consumer<SyncProvider>(
      builder: (context, syncProvider, _) {
        if (!syncProvider.isInitialized) {
          return const SizedBox.shrink();
        }

        if (syncProvider.isSyncing) {
          return const _SyncInProgressIndicator();
        }

        if (syncProvider.hasPendingOperations) {
          return const _PendingSyncIndicator();
        }

        return const _SyncedIndicator();
      },
    );
  }
}

/// Indicator shown when sync is in progress
class _SyncInProgressIndicator extends StatelessWidget {
  const _SyncInProgressIndicator();

  @override
  Widget build(BuildContext context) {
    return Tooltip(
      message: 'Synchronizing with server...',
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          const SizedBox(
            width: 16,
            height: 16,
            child: CircularProgressIndicator(
              strokeWidth: 2,
            ),
          ),
          const SizedBox(width: 8),
          Text(
            'Syncing...',
            style: Theme.of(context).textTheme.bodySmall,
          ),
        ],
      ),
    );
  }
}

/// Indicator shown when there are pending operations
class _PendingSyncIndicator extends StatelessWidget {
  const _PendingSyncIndicator();

  @override
  Widget build(BuildContext context) {
    return Tooltip(
      message: 'Changes pending synchronization',
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            Icons.sync_problem,
            size: 16,
            color: Theme.of(context).colorScheme.warning ?? Colors.orange,
          ),
          const SizedBox(width: 8),
          Text(
            'Pending sync',
            style: Theme.of(context).textTheme.bodySmall?.copyWith(
                  color: Theme.of(context).colorScheme.warning ?? Colors.orange,
                ),
          ),
        ],
      ),
    );
  }
}

/// Indicator shown when all data is synced
class _SyncedIndicator extends StatelessWidget {
  const _SyncedIndicator();

  @override
  Widget build(BuildContext context) {
    final syncProvider = Provider.of<SyncProvider>(context, listen: false);
    final lastSync = syncProvider.lastSyncTime;
    
    String syncTimeText = 'Never synced';
    if (lastSync != null) {
      final now = DateTime.now();
      final difference = now.difference(lastSync);
      
      if (difference.inMinutes < 1) {
        syncTimeText = 'Just now';
      } else if (difference.inHours < 1) {
        syncTimeText = '${difference.inMinutes}m ago';
      } else if (difference.inDays < 1) {
        syncTimeText = '${difference.inHours}h ago';
      } else {
        syncTimeText = '${difference.inDays}d ago';
      }
    }
    
    return Tooltip(
      message: 'All changes synchronized with server',
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            Icons.sync,
            size: 16,
            color: Theme.of(context).colorScheme.primary,
          ),
          const SizedBox(width: 8),
          Text(
            'Synced $syncTimeText',
            style: Theme.of(context).textTheme.bodySmall?.copyWith(
                  color: Theme.of(context).colorScheme.primary,
                ),
          ),
        ],
      ),
    );
  }
}

// Extension to add warning color to ColorScheme
extension ColorSchemeExtension on ColorScheme {
  Color? get warning => Colors.orange;
}