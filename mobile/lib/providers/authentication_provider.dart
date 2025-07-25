import 'package:flutter/material.dart';
import '../repositories/subscription_repository.dart';
import '../repositories/label_repository.dart';
import '../repositories/family_repository.dart';
import '../services/authentication_service.dart';

class AuthenticationProvider extends ChangeNotifier {
  AuthenticatedUser? _user;
  bool _isAuthenticated = false;
  final AuthenticationService _authenticationService;
  final SubscriptionRepository? _subscriptionRepository;
  final LabelRepository? _labelRepository;
  final FamilyRepository? _familyRepository;

  AuthenticationProvider({
    required AuthenticationService authenticationService,
    SubscriptionRepository? subscriptionRepository,
    LabelRepository? labelRepository,
    FamilyRepository? familyRepository,
  })  : _authenticationService = authenticationService,
        _subscriptionRepository = subscriptionRepository,
        _labelRepository = labelRepository,
        _familyRepository = familyRepository {
    refreshUser();
  }

  AuthenticatedUser? get user => _user;

  bool get isAuthenticated => _isAuthenticated;

  Future<void> refreshUser() async {
    final isAuthenticated = await _authenticationService.isAuthenticated();

    if (isAuthenticated) {
      _user = await _authenticationService.getCurrentUser();
      _isAuthenticated = true;

      // Switch to user-specific data
      if (_user != null) {
        await _switchToUserData(_user!.id);
      }
    } else {
      _user = null;
      _isAuthenticated = false;

      // Switch to anonymous data
      await _switchToUserData(null);
    }

    notifyListeners();
  }

  Future<void> signIn() async {
    final token = await _authenticationService.login();
    if (token != null) {
      await refreshUser();
    }
  }

  Future<void> signOut() async {
    await _authenticationService.logout();

    // Clear user data before switching to anonymous mode
    await _clearUserData();

    _user = null;
    _isAuthenticated = false;

    // Switch to anonymous data
    await _switchToUserData(null);

    notifyListeners();
  }

  // Switch to user-specific data or anonymous data
  Future<void> _switchToUserData(String? userId) async {
    // Update repositories to use user-specific boxes
    if (_subscriptionRepository != null) {
      await _subscriptionRepository!.setCurrentUser(userId);
    }

    if (_labelRepository != null) {
      await _labelRepository!.setCurrentUser(userId);
    }

    if (_familyRepository != null) {
      await _familyRepository!.setCurrentUser(userId);
    }
  }

  // Clear user data when logging out
  Future<void> _clearUserData() async {
    if (_user == null) return;

    // Clear data in repositories
    if (_subscriptionRepository != null) {
      await _subscriptionRepository!.clearUserData();
    }

    if (_labelRepository != null) {
      await _labelRepository!.clearUserData();
    }

    if (_familyRepository != null) {
      await _familyRepository!.clearUserData();
    }
  }
}
