import 'package:kinde_flutter_sdk/kinde_flutter_sdk.dart';

class AuthenticatedUser {
  final String id;
  final String email;
  final String? displayName;

  AuthenticatedUser({required this.id, required this.email, this.displayName});
}

class AuthenticationService {
  final KindeFlutterSDK sdk = KindeFlutterSDK.instance;

  Future<bool> isAuthenticated() async {
    return await sdk.isAuthenticated();
  }

  Future<AuthenticatedUser?> getCurrentUser() async {
    try {
      final isAuthenticated = await sdk.isAuthenticated();

      if (isAuthenticated) {
        final kindeUser = sdk.getUserDetails();

        if (kindeUser != null) {
          return AuthenticatedUser(
            id: kindeUser.id ?? '',
            email: kindeUser.email ?? '',
            displayName: kindeUser.givenName,
          );
        }
      }
      return null;
    } catch (e) {
      return null;
    }
  }

  Future<String?> getToken() async {
    return await sdk.getToken();
  }

  Future<void> deleteAccount() async {
    // todo
  }

  Future<String?> login() async {
    try {
      return await sdk.login(type: AuthFlowType.pkce);
    } catch (e) {
      rethrow;
    }
  }

  Future<void> logout() async {
    await sdk.logout();
  }
}
