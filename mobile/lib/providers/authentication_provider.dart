import 'package:flutter/material.dart';
import '../services/authentication_service.dart';

class AuthenticationProvider extends ChangeNotifier {
  AuthenticatedUser? _user;
  bool _isAuthenticated = false;
  final AuthenticationService _authenticationService;

  AuthenticationProvider({required AuthenticationService authenticationService})
    : _authenticationService = authenticationService {
    refreshUser();
  }

  AuthenticatedUser? get user => _user;

  bool get isAuthenticated => _isAuthenticated;

  Future<void> refreshUser() async {
    final isAuthenticated = await _authenticationService.isAuthenticated();

    if (isAuthenticated) {
      _user = await _authenticationService.getCurrentUser();
      _isAuthenticated = true;
    } else {
      _user = null;
      _isAuthenticated = false;
    }

    notifyListeners();
  }

  // Sign out the user
  Future<void> signOut() async {
    await _authenticationService.logout();
    _user = null;
    _isAuthenticated = false;
    notifyListeners();
  }
}
