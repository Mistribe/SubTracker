// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
/// Retrieve a paginated list of labels with optional filtering by owner type
class LabelsRequestBuilderGetQueryParameters implements AbstractQueryParameters {
    ///  Family ID (UUID format)
    String? familyId;
    ///  Number of items (default: 10)
    int? limit;
    ///  Offset (default: 0)
    int? offset;
    ///  Owner types to filter by (system,personal,family). Can be provided multiple times.
    /// @QueryParameter('owner_type')
    List<String>? ownerType;
    /// Extracts the query parameters into a map for the URI template parsing.
    @override
    Map<String, dynamic> toMap() {
        return {
            'familyId' : familyId,
            'limit' : limit,
            'offset' : offset,
            'owner_type' : ownerType,
        };
    }
}
