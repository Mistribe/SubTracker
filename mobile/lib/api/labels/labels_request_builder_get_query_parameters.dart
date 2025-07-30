// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
/// Get all labels
class LabelsRequestBuilderGetQueryParameters implements AbstractQueryParameters {
    ///  Offset number
    int? page;
    ///  Number of items per page
    int? size;
    ///  Include default labels
    /// @QueryParameter('with_default')
    bool? withDefault;
    /// Extracts the query parameters into a map for the URI template parsing.
    @override
    Map<String, dynamic> toMap() {
        return {
            'page' : page,
            'size' : size,
            'with_default' : withDefault,
        };
    }
}
