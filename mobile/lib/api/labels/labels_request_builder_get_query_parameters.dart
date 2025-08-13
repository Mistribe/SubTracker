// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
/// Retrieve a paginated list of labels with optional filtering by owner type and search text
class LabelsRequestBuilderGetQueryParameters implements AbstractQueryParameters {
    ///  Maximum number of items to return (default: 10)
    int? limit;
    ///  Number of items to skip for pagination (default: 0)
    int? offset;
    ///  Search text to filter labels by name
    String? search;
    /// Extracts the query parameters into a map for the URI template parsing.
    @override
    Map<String, dynamic> toMap() {
        return {
            'limit' : limit,
            'offset' : offset,
            'search' : search,
        };
    }
}
