# Localize
Utility to convert csv data into android, ios or web string resource formats.

## Motivation
A lot of the times we build applications that spans across countries and regions, and we want to add local languages (i.e localization) translation in our application. Nowadays, there are tools to manage localization, which easily integrate with CI CD pipelines to automatically update content, even so they sometimes provide resource files for platform specific format. Sometimes, they dont, and one of the most common formats you can get data is csv. Hence, this tool is written focusing csv file as input, which is described in detail later. This tool can be used to automate updating the string resources from same source, without manually updating the file each time. 

## Sneak Peek
![demo](https://github.com/talhahasanzia/string-resource-utility/blob/main/assets/screen_shots.gif)

## Build
- Install [Go](https://go.dev). 
- Goto working directory.
- To run the utility `go run localize.go -file=strings.csv -platform=ios`
- To build an executable, `go build localize` and you can redistribute.

## Usage
- Prepare a csv file with 1st column as key, and subsequent columns as locales e.g. Note: **key** and **locale** will be trimmed.

| key           | en           | es  | it  |
| ------------- |------------- | ----- | ----- |
| welcome_text  | welcome      | bienvenidos | benvenuto |
| bye_text      | goodbye      |   adiós | arrivederci |

- run `./localize -file=<PATH_TO_CSV> -platform=<IOS_ANDROID_WEB> -debug=<TRUE_FALSE>`
- `-file`: Absolute path to csv file
- `-output`: output directory for generated files
- `-overwrite`: overwrite existing file contents: true/false  
- `-platform`: platform for which string resource will be generated. Options: `ios`,`android` and `web` (anyone at a time).
- `-debug`: Show logs while creating resources, by default it is `false`. Options: `true` or `false`
- Currently, it will generate resource files in the same folder where the `localize` command is executed.



## Resource files
### Android
If you chose `-platfrom=android`, the generated files will be `strings.xml`, since we are targeting multiple locales, this tool will generate `strings_en.xml` for `en` locale, and so on. Here is what this generated file looks like when above mentioned csv is used:
```
<resources>
<string name="welcome_text">welcome</string>
<string name="bye_text">goodbye</string>
</resources>
```

Specifically for android, you have to specify keys without spaces.

### iOS
Similarly, if you chose `-platfrom=ios`, the generated files will be `Localized.strings`, since we are targeting multiple locales, this tool will generate `Localized_en.strings` for `en` locale, and so on. Here is what this generated file looks like when above mentioned csv is used:
```
"welcome_text" = "welcome";
"bye_text" = "goodbye";
```

### Flutter (arb files only)
Similarly, if you chose `-platfrom=flutter`, the generated files will be `app_<en>.arb`, since we are targeting multiple locales, this tool will generate `app_en.arb` for `en` locale, and so on. Here is what this generated file looks like when above mentioned csv is used:
```
{
"welcome_text": "welcome",
"bye_text": "goodbye",
"eof":"eof"
}
```

Note that a dummy entry as "eof" is set in the end of arb file in order to mitigate error resulting from last left comma (*,*) when `flutter gen-l10n` is run.

### Web
And finally, if you chose `-platfrom=web` (assuming you are working on a typescript framework), the generated files will be `strings.ts`, since we are targeting multiple locales, this tool will generate `strings_en.ts` for `en` locale, and so on. Here is what this generated file looks like when above mentioned csv is used:
```
const LOCALIZED_STRINGS = {
welcome_text: "welcome",
bye_text: "goodbye",
}

export default LOCALIZED_STRINGS;
```
## Customization
You can always customize the way this tool work, following are the files where you will find relevant logic:

- `reader.go`: Reads data from `.csv` file and returns a 2d array.
- `parser.go`: Gets 2d array data as input and returns a list of `Record` object
- `writer.go`: Writes to file for each `Record` entry, and also responsible for opening and closing data specific to platform.

### Explanation
- Once the data is ready to be written in the form of list (slice) of `Record`, the code groups records by locale.
- Each locale's records are processed sequentially, writing all entries for one locale before moving to the next.
- For each locale, the code creates platform + locale specific files and populates some pretext if needed by that specific platform.
- The writer processes all records for the current locale in a simple loop, writing each entry to the appropriate file.
- This is achieved by grouping records using a map where locale is the key and a slice of records is the value. So `localeRecords[locale]` contains all entries for that locale.
- After all records for a locale are written, `CloseFile` is called which writes the closing statement (if needed by platform) to the platform + locale specific file.
- This sequential approach is simpler to understand and debug while producing the same output as the previous concurrent implementation.
- Sequential approach also fixes an issue when sometimes, last entry is placed outside closing token.

## Releases
Find releases for your specific platforms in [releases](https://github.com/talhahasanzia/string-resource-utility/releases/tag/v1.0) section.
`.zip` file contains executables for 64-bit (marked as `amd64`) and 32-bit (marked as `386`) versions. `arm` is target release for Apple M1 Macbooks. `.exe` for Microsoft Windows.

## Get Involved
Feel free to create PRs and report issues or submit feedback. Just remember to be polite and respectful :)
