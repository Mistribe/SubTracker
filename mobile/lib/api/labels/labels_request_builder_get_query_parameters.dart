// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
/// Retrieve a paginated list of labels with optional filtering by owner type
class LabelsRequestBuilderGetQueryParameters implements AbstractQueryParameters {
    ///  Owner types to filter by (system,personal,family). Can be provided multiple times.
    /// @QueryParameter('owner_type')
    List<String>? ownerType;
    ///  Page number (default: 1)
    int? page;
    ///  Number of items per page (default: 10)
    int? size;
    /// Extracts the query parameters into a map for the URI template parsing.
    @override
    Map<String, dynamic> toMap() {
        return {
            'owner_type' : ownerType,
            'page' : page,
            'size' : size,
        };
    }
}
