// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
/// Get information about a family invitation using invitation code
class InvitationRequestBuilderGetQueryParameters implements AbstractQueryParameters {
    ///  Invitation code
    String? code;
    ///  Family member ID
    /// @QueryParameter('family_member_id')
    String? familyMemberId;
    /// Extracts the query parameters into a map for the URI template parsing.
    @override
    Map<String, dynamic> toMap() {
        return {
            'code' : code,
            'family_member_id' : familyMemberId,
        };
    }
}
