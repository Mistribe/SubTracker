import 'package:flutter/material.dart';
import 'package:kinde_flutter_sdk/kinde_flutter_sdk.dart';
import '../models/user.dart';

class User {
  final String id;
  final String email;
  final String? displayName;

  User({required this.id, required this.email, this.displayName});
}

class UserProvider extends ChangeNotifier {
  User? _user;
  bool _isAuthenticated = false;

  UserProvider() {
    // Initialize by checking if user is already authenticated
    refreshUser();
  }

  // Get the current user
  User? get user => _user;

  // Check if a user is authenticated
  bool get isAuthenticated => _isAuthenticated;

  // Refresh the user state from Kinde
  Future<void> refreshUser() async {
    try {
      final isAuthenticated = await KindeFlutterSDK.instance.isAuthenticated();

      if (isAuthenticated) {
        final kindeUser = KindeFlutterSDK.instance.getUserDetails();

        if (kindeUser != null) {
          _user = User(
            id: kindeUser.id ?? '',
            email: kindeUser.email ?? '',
            displayName: kindeUser.givenName,
          );
          _isAuthenticated = true;
        } else {
          _user = null;
          _isAuthenticated = false;
        }
      } else {
        _user = null;
        _isAuthenticated = false;
      }

      notifyListeners();
    } catch (e) {
      _user = null;
      _isAuthenticated = false;
      notifyListeners();
    }
  }

  // Sign out the user
  Future<void> signOut() async {
    try {
      await KindeFlutterSDK.instance.logout();
      _user = null;
      _isAuthenticated = false;
      notifyListeners();
    } catch (e) {
      // Handle error
      print('Error signing out: $e');
    }
  }
}
