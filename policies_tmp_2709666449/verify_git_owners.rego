package policies.verify_git_owners

import data
import future.keywords.in
import future.keywords.every
import data.policies.sbom_parser as parser

default allow = false

list = l {
    branch := checkInputExists(input.match, "git_branch")
    tag := checkInputExists(input.match, "git_tag")
    commit := checkInputExists(input.match, "git_commit")
    url := checkInputExists(input.match, "git_url")
    format := checkInputExists(input.match, "format")
    
    l := {
        "target_git_url": url,        
        "target_git_branch": branch,
        "target_git_tag": tag,
        "target_git_commit": commit,
        "content_type": format,
        "predicate_type": "https://cyclonedx.org/bom",
        "input_scheme": "git"
    }
}

allow {  
    rule_input := input.rule_input
    statement := input.rule_evidences[0]    

    bom := parser.getBom(statement)
    parser.validateBom(bom) 
    parser.checkGitScheme(bom)

    # verify git signatures
    notSigned := parser.getNotSigneddFiles(bom)
    notSignedCount := count(notSigned)
    any([notSignedCount == 0, rule_input.signed_commits == false]) # if signed_commits is true, check the notSignedCount, otherwise, ignore it

    # verify git owners
    unallowedOwners := checkOwners(rule_input, bom)  
    count(unallowedOwners) == 0    
}


# debug = d {
#     statement := input.rule_evidences[0]    
#     bom := parser.getBom(statement)
#     rule_input := input.rule_input
#     adminOwners := getAdmins(rule_input)
#     defualtOwners := getDefaultOwners(rule_input)
#     specificOwners := getSpecificOwners(rule_input)
#     files := parser.getFilesLastCommit(bom)
#     checkedOwners := checkOwners(rule_input, bom)    

#     d := {
#         "admins": adminOwners,
#         "default": defualtOwners,
#         "countCheckedOwners": count(checkedOwners),
#         "countDefault": count(defualtOwners),
#         "countSpec": count(specificOwners),
        
#         "checkedOwners": checkedOwners,
#     }
# }

results_message_scheme_error = v {
    statement := input.rule_evidences[0]
    bom := parser.getBom(statement)
    not parser.checkGitScheme(bom)  
    scheme := parser.getScheme(bom)
    v := sprintf("scheme is %v instead of %v",[scheme, parser.gitScheme])
}

results_message_sbom_error = v {
    statement := input.rule_evidences[0]
    bom := parser.getBom(statement)
    not parser.validateBom(bom)
    v := "input is not a SBOM"
}

######################################
# verify git signatures
######################################
results_message_verify_signatures = v {
    rule_input := input.rule_input
    statement := input.rule_evidences[0]
    bom := parser.getBom(statement)    
    rule_input.signed_commits
    notSigned := parser.getNotSigneddFiles(bom)
    not count(notSigned) == 0    
    files := parser.getFiles(bom)    
    v :=  sprintf("%v/%v files were commited by unsigned commits.",[count(notSigned), count(files)])
}

results_details_verify_signatures = d {
    rule_input := input.rule_input
    statement := input.rule_evidences[0]
    bom := parser.getBom(statement)
    rule_input.signed_commits
    notSigned := parser.getNotSigneddFiles(bom)
    not count(notSigned) == 0
    d := {
        "unsigned_files": parser.getNotSignedFilesInfo(bom),
    }
}

######################################

######################################
# verify git owners
######################################

results_message_verify_owners = v {
    rule_input := input.rule_input
    statement := input.rule_evidences[0]
    bom := parser.getBom(statement)
    files := parser.getFilesLastCommit(bom)
    unallowedOwners := checkOwners(rule_input, bom)  
    v := sprintf("%v/%v files were commited by unauthorized owners.",[count(unallowedOwners), count(files)])
}



results_details_verify_owners = d {
    rule_input := input.rule_input
    statement := input.rule_evidences[0]
    bom := parser.getBom(statement)
    files := parser.getFilesLastCommit(bom)
    unallowedOwners := checkOwners(rule_input, bom)  
    d := {
        "unallowed_files": unallowedOwners,
    }
}
######################################

######################################
# collect results
######################################
results_message = v {
    v := [
        results_message_scheme_error,
        results_message_sbom_error,     
    ]
}

results_message = v {
    not results_message_scheme_error
    not results_message_sbom_error
    v := [
        results_message_verify_signatures,
        results_message_verify_owners,
    ]
}

results_details = d {
    d := [
        results_details_verify_signatures,
        results_details_verify_owners,
    ]
}

######################################

######################################
# functions
######################################

getAdmins(rule_input) = admins {
    admins := checkInputExists(rule_input, "admin")     
}

getDefaultOwners(rule_input) = def {
    def := checkInputExists(rule_input, "default")     
}

getSpecificOwners(rule_input) = owners {
    owners := checkInputExists(rule_input, "specific")     
}

collectInfoForCommits(statement, commits) = info {
    bom := parser.getBom(statement)
    info := { x |
        some i
        commitHash := commits[i]
        files := parser.getFilesForGitCommit(bom, commitHash)
        author := parser.getCommitAuthor(bom, commitHash) 
        x := {            
            "commit": commitHash,
            "author": author,
            "files": files,
        }
    }
}

getOwnersTotal(bom) = o {    
    owners := parser.getLastCommitOwners(bom)
    o := owners
}


checkSignedCommits(bom) {
    notSigned := parser.getNotSigneddFiles(bom)
    count(notSigned) == 0
}

defined(path, field) {
    _ = path[field] 
}

checkInputExists(path, field) = v {
	path[field]
    v := path[field]
}

checkInputExists(path, field) = v {
	not defined(path, field)
    v := ""
}


checkSpecificOwners(rule_input, fileEntry) = o {
    specificOwners := getSpecificOwners(rule_input)
    count(specificOwners) > 0 # ignore default owners
    o := { x |
        some i,j
        ownerEntry := specificOwners[i]
        ownersList := ownerEntry["owners"]
        startswith(fileEntry.file, ownerEntry["path"])
        ownersList[j] == fileEntry["owner"]
        x := fileEntry
    }
}

checkDefaultOwners(rule_input, fileEntry) = o {
    defualtOwners := getDefaultOwners(rule_input)
    o := { x |
        some i
        owner := defualtOwners[i]        
        owner == fileEntry["owner"]
        x := fileEntry
    }
}

checkAdminOwners(rule_input, fileEntry) = o {
    adminOwners := getAdmins(rule_input)
    o := { x |
        some i
        owner := adminOwners[i]        
        owner == fileEntry["owner"]
        x := fileEntry
    }
}

checkSpecificOrDefault(rule_input, fileEntry) = o {
    # check specific owners if more than one is defined
    specificOwners := getSpecificOwners(rule_input)
    count(specificOwners) > 0
    o := checkSpecificOwners(rule_input, fileEntry)
}

checkSpecificOrDefault(rule_input, fileEntry) = o {
    # check default owners if no specific owners were defined
    specificOwners := getSpecificOwners(rule_input)
    count(specificOwners) == 0
    o := checkDefaultOwners(rule_input, fileEntry)
}

checkOwners(rule_input, bom) = c {
    files := parser.getFilesLastCommit(bom)
    #check specific or default owners against the files
    verified_files := { x[_] |
        some i
        fileEntry := files[i]
        x := checkSpecificOrDefault(rule_input, fileEntry)
    }
    rest := files - verified_files

    #check admins against the rest of the files
    verified_admins := { x[_] |
        some i
        fileEntry := rest[i]
        x := checkAdminOwners(rule_input, fileEntry)
    }
    c := rest - verified_admins    
}