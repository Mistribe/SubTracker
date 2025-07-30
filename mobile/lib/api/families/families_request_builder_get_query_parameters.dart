// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
/// Retrieve a paginated list of families for the authenticated user
class FamiliesRequestBuilderGetQueryParameters implements AbstractQueryParameters {
    ///  Number of items per page (default: 10)
    int? limit;
    ///  Page number (default: 1)
    int? offset;
    /// Extracts the query parameters into a map for the URI template parsing.
    @override
    Map<String, dynamic> toMap() {
        return {
            'limit' : limit,
            'offset' : offset,
        };
    }
}
