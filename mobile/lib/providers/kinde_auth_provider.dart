import 'package:microsoft_kiota_bundle/microsoft_kiota_bundle.dart';

import '../services/authentication_service.dart';

class KindeAuthenticationProvider implements AuthenticationProvider {
  final AuthenticationService authenticationService;

  KindeAuthenticationProvider(this.authenticationService);

  @override
  Future<void> authenticateRequest(
    RequestInformation request, [
    Map<String, Object>? additionalAuthenticationContext,
  ]) async {
    final token = await authenticationService.getToken();
    if (token != null) {
      request.headers["Authorization"] = {"Bearer $token"};
    }
  }
}
