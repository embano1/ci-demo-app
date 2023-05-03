package policies.sbom_parser

# import data
import future.keywords.in
import future.keywords.every

# default allow = false

format = "CycloneDX"
versions = {"1.0","1.1","1.2","1.3", "1.4"}
gitScheme = "git"
dirScheme = "dir"
supportedTools = {"gensbom","valint"}
typeCommit = "commit"
typeFile = "file"
typeLibrary = "library"
pkgNpm = "npm"
#intoto-statement support
statement_type = "https://in-toto.io/Statement/v0.1"
bom_predicate = "https://cyclonedx.org/bom"


checkIntoto(statement) {
    statement._type == statement_type
}

checkBomPredicate(statement) {
    statement.predicateType == bom_predicate
}


checkBomStatement(statement) {
    checkIntoto(statement)
    checkBomPredicate(statement)
}

getBom(statement) = bom {
    bom := statement.predicate.bom
}

getBom(statement) = statement {
    not statement.predicate.bom
}

getEnv(statement) = env {
    env := statement.predicate.environment
}

checkVersion(specVersion) {
  versions[specVersion]
}

checkFormat(f) {
  format == f
}

checkScheme(bom)  {    
    scheme := getScheme(bom)
    scheme == gitScheme
}

checkTool(bom) {
    tool := getToolName(bom)
    supportedTools[tool]
}

getScheme(bom) = scheme {
    scheme := getGitMetadataPropertyByName(bom, "input_scheme")
}

getMetadataProperties(bom) = props {            
    metaComp := getMetadataComponent(bom)
    props := metaComp.properties
}

getMetadata(bom) = metadata {
    metadata := bom.metadata
}

getComponents(bom) = comps {
    comps := bom.components
}

getMetadataComponent(bom) = metadata_component {
        metadata := getMetadata(bom)        
        metadata_component := metadata.component
}

getSubject(bom) = subs {
    comp := getMetadataComponent(bom)
    name := comp["name"]
    version := split(comp["version"], ":")    
    subs := {
        "name": name,
        "alg": version[0],
        "hash": version[1]
    }
}

getContentType(statement) = contentType {
    env := getEnv(statement)
    contentType := env.content_type
}

getToolName(bom) = toolName {       
    metadata := getMetadata(bom)        
    toolName := metadata.tools[_].name
}

getTimestamp(bom) = timestamp {
    metadata := getMetadata(bom)        
    timestamp := metadata.timestamp
}

validateBom(bom) {
    checkVersion(bom.specVersion)
    checkFormat(bom.bomFormat)
    checkScheme(bom)
    checkTool(bom)
}

checkGitScheme(bom) {
    scheme := getScheme(bom)
    scheme == gitScheme
}


getGitTarget(bom) = gitUrl {
    checkGitScheme(bom)
    metaComp := getMetadataComponent(bom)
    gitUrl := metaComp.name
}

getGitBranch(bom) = gitBranch {
    gitBranch := getGitMetadataPropertyByName(bom, "git_branch")
}

getGitTag(bom) = gitTag {
    gitTag := getGitMetadataPropertyByName(bom, "git_tag")
}

getGitTag(bom) = "" { #return empty string if git_tag doesn't exist
    not getGitMetadataPropertyByName(bom, "git_tag")
}

getGitCommit(bom) = gitCommit {
    gitCommit := getGitMetadataPropertyByName(bom, "git_commit")
}

getGitMetadataPropertyByName(bom, propertyName) = propertyValue {
    props := getMetadataProperties(bom)    
    x := props[i]
    x.name == propertyName     
    propertyValue := x.value   
}


getComponentsByType(bom, type) = componentMap {
componentMap := { x["bom-ref"]: x | 
        some i
        x := bom.components[i]
        x.type == type 
    }
}

getComponentsByGroup(bom, group) = componentMap {
    componentMap := { x["bom-ref"]: x | 
        some i
        x := bom.components[i]
        x.group == group 
    }
}

getComponentsByTypeAndGroup(bom, type, group) = componentMap {
    componentMap := { x["bom-ref"]: x | 
        some i
        x := bom.components[i]
        x.group == group
        x.type == type 
    }
}

getGitCommits(bom) = commits {
     commits := getComponentsByTypeAndGroup(bom, typeCommit, typeCommit)
}

getFiles(bom) = files {
     files := getComponentsByTypeAndGroup(bom, typeFile, typeFile)
}

getGitCommitHashes(bom) = commitHashes {
    commitHashes := {x["name"]: x["properties"] |
        some i
        commits := getGitCommits(bom)
        x := commits[i]
    }
}

getGitFilesLastCommit(bom) = filesLastCommit {    
    filesLastCommit := { x["name"]: lastCommit |
        some i,j
        files := getComponentsByTypeAndGroup(bom, typeFile, typeFile)
        x := files[i]
        prop := x.properties
        prop[j].name == "last_commit"
        lastCommit := prop[j].value
    }
}

getLastCommitOwners(bom) = owners {    
    filesLastCommit := getGitFilesLastCommit(bom)  
    owners := { x |
        some i
        fileLastCommitHash := filesLastCommit[i]
        # files := getFilesForGitCommit(bom, fileLastCommitHash)  
        # author := getCommitAuthor(bom, fileLastCommitHash)
        email :=  getCommitAuthorEmail(bom, fileLastCommitHash)
        x := {
            # "files": files,
            "commit": fileLastCommitHash,
            "email": email            
        }
    }
}

getFilesLastCommit(bom) = files { 
    filesLastCommit := getGitFilesLastCommit(bom)
    files := { x |
        some i,j
        fileLastCommitHash := filesLastCommit[i]
        fs := getFilesForGitCommit(bom, fileLastCommitHash)
        email :=  getCommitAuthorEmail(bom, fileLastCommitHash)
        x := {
            "file": fs[j],
            "commit": fileLastCommitHash,
            "owner": email
        } 
    }
}



getNotSigneddFiles(bom) = notSigneddFiles {    
    notSigneddFiles :=  { x[_] |
        filesLastCommit := getGitFilesLastCommit(bom)        
        some i
        fileLastCommitHash := filesLastCommit[i]
        not checkIfCommitIsSigned(bom, fileLastCommitHash)
        x := getFilesForGitCommit(bom, fileLastCommitHash)
    }
}

getNotSignedFilesInfo(bom) = notSigneddFilesInfo {
    notSigneddFiles := getNotSigneddFiles(bom)
    filesLastCommit := getGitFilesLastCommit(bom)     
    notSigneddFilesInfo := { x |
        some i
        file := notSigneddFiles[i]        
        commitHash := filesLastCommit[file]
        author := getCommitAuthor(bom, commitHash)
        x := {        
            "file": file,
            "commit": commitHash,
            "author": author
        }
    }
}



getFilesForGitCommit(bom, commitHash) = files {
    files := { x |
        filesLastCommit := getGitFilesLastCommit(bom)
        some file; filesLastCommit[file] == commitHash
        x := file
    }
}

checkIfCommitIsSigned(bom, commitHash) {
    sig := getCommitSignature(bom, commitHash)
    not sig == ""    
}

getCommitSignature(bom, commitHash) = sig {
    commitHashes := getGitCommitHashes(bom)
    some i
    commitHashes[commitHash][i].name == "PGPSignature"
    sig := commitHashes[commitHash][i].value
}

getCommitAuthor(bom, commitHash) = author {
    commitHashes := getGitCommitHashes(bom)
    some i
    commitHashes[commitHash][i].name == "Author"
    unparsedAuthor := commitHashes[commitHash][i].value    
    authorName := split(unparsedAuthor, " \u003c")
    authorEmail := replace(authorName[1],"\u003e","")
    author := {
        "name": authorName[0],
        "e-mail":authorEmail
    }
}

getCommitAuthorEmail(bom, commitHash) = email {
    commitHashes := getGitCommitHashes(bom)
    some i
    commitHashes[commitHash][i].name == "Author"
    unparsedAuthor := commitHashes[commitHash][i].value    
    authorName := split(unparsedAuthor, " \u003c")
    email := replace(authorName[1],"\u003e","")
}

getNpmPackages(bom) = npmPackages {
    packages :=  getComponentsByTypeAndGroup(bom, typeLibrary, pkgNpm)   
    npmPackages := { x |
        some i        
        imp_path := getPackagePropertyByName(packages[i],"importer-path")
        x := {
            "name": packages[i].name,
            "version": packages[i].version,
            "importer-path": imp_path,
        }
    }

}

getBomPackages(bom) = allPackages {
    packages :=  getComponentsByType(bom, typeLibrary)   
    allPackages := { x |
        some i        
        imp_path := getPackagePropertyByName(packages[i],"importer-path")
        x := {
            "name": packages[i].name,
            "version": packages[i].version,
            "group": packages[i].group,
            "importer-path": imp_path,
        }
    }
}

getPackagePropertyByName(pkg, propertyName) = propertyValue {
    some i
    x := pkg.properties[i]
    x.name == propertyName     
    propertyValue := x.value   
}

getMismatchedPackages(src, dst) = mismatchList {    
    mismatchList := src - dst
}

getMatchedPackages(src, dst) = matchedList {
    mismatched := getMismatchedPackages(src, dst)
    matchedList := src - mismatched
}
